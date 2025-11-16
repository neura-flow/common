package election

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/util"
	"github.com/samuel/go-zookeeper/zk"
)

type Election struct {
	ctx        context.Context
	config     *Config
	logger     log.Logger
	lock       *ResourceLock
	resourceId string //resourceId 保存到 zk node 中,用于判断当前的连接是否选主成功的连接
	isLeader   bool
}

func NewElection(ctx context.Context, logger log.Logger, cfg *Config) (*Election, error) {
	if err := validate(cfg); err != nil {
		return nil, err
	}
	lock, err := NewResourceLock(logger, cfg.ZkServers)
	if err != nil {
		return nil, err
	}
	return &Election{
		ctx:        ctx,
		logger:     logger,
		config:     cfg,
		lock:       lock,
		resourceId: util.GUID(),
	}, nil
}

func (le *Election) Run() {
	for {
		select {
		case <-le.ctx.Done():
			le.release()
			return
		default:
			le.run(le.ctx)
		}
	}
}

func (le *Election) run(ctx context.Context) {
	defer func() {
		le.isLeader = false
		le.config.Callbacks.OnStoppedLeading()
	}()

	acquired, err := le.acquire()
	if err != nil || !acquired {
		time.Sleep(time.Duration(2) * time.Second)
		return
	}

	le.isLeader = true
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	go le.config.Callbacks.OnStartedLeading(ctx)

	err = le.watch(ctx)
	if err != nil {
		le.logger.Debugf("failed to renew %s leader lease,err: %v", le.config.Identity, err)
		time.Sleep(time.Duration(2) * time.Second)
		return
	}
}

func validate(cfg *Config) error {
	if cfg == nil {
		return errors.New("cfg is required")
	}
	if util.IsBlank(&cfg.ZkServers) {
		return errors.New("zk servers is required")
	}
	if util.IsBlank(&cfg.ElectionRoot) {
		return errors.New("root path is required")
	}
	if !strings.HasPrefix(cfg.ElectionRoot, "/") {
		return errors.New("root path should begin with '/'")
	}
	if util.IsBlank(&cfg.ElectionID) {
		return errors.New("leaderElectionID is required")
	}
	return nil
}

func (le *Election) IsLeader() bool {
	return le.isLeader
}

func (le *Election) acquire() (bool, error) {
	if err := le.ensureRoot(); err != nil {
		return false, err
	}
	fullElectionID := le.getFullElectionID()
	if err := le.elected(fullElectionID); err == nil {
		return true, nil
	} else if errors.Is(err, zk.ErrNoNode) {
		created, err := le.lock.Create(fullElectionID, []byte(le.resourceId), FlagEphemeral)
		if err != nil {
			return false, err
		}
		if !strings.EqualFold(created, fullElectionID) {
			return false, fmt.Errorf("created node mismatch, want: %s created: %s", fullElectionID, created)
		}
	}
	if err := le.elected(fullElectionID); err != nil {
		return false, err
	}
	le.logger.Infof("%s acquired the leader lease", le.config.Identity)
	return true, nil
}

func (le *Election) elected(fullElectionID string) error {
	if data, err := le.lock.Get(fullElectionID); err != nil {
		return err
	} else if !strings.EqualFold(string(data), le.resourceId) {
		return fmt.Errorf("failed to acquire %s", fullElectionID)
	}
	return nil
}

func (le *Election) watch(ctx context.Context) error {
	childCh, err := le.lock.Watch(le.getFullElectionID())
	if err != nil {
		return err
	}
	for {
		select {
		case event := <-childCh:
			if event.Type == zk.EventNodeDeleted {
				return errors.New(event.Type.String())
			} else if event.State != zk.StateConnected && event.State != zk.StateHasSession {
				return errors.New(event.State.String())
			} else {
				continue
			}
		case <-ctx.Done():
			le.logger.Infof("%s receive cancel", le.config.Identity)
			return nil
		}
	}
}

func (le *Election) release() {
	if le.lock != nil {
		le.lock.Close()
	}
}

func (le *Election) ensureRoot() error {
	if exists, err := le.lock.Exists(le.config.ElectionRoot); err != nil {
		return err
	} else if !exists {
		created, err := le.lock.Create(le.config.ElectionRoot, nil, FlagPermanent)
		if err != nil {
			return err
		}
		if !strings.EqualFold(created, le.config.ElectionRoot) {
			return fmt.Errorf("failed to created root node, identity: %s want: %s created: %s", le.config.Identity, le.config.ElectionRoot, created)
		}
	}
	return nil
}

func (le *Election) getFullElectionID() string {
	return fmt.Sprintf("%s/%s", le.config.ElectionRoot, le.config.ElectionID)
}

type Config struct {
	ZkServers    string
	ElectionRoot string
	ElectionID   string
	Callbacks    Callbacks
	Identity     string
}

type Callbacks struct {
	OnStartedLeading func(context.Context)
	OnStoppedLeading func()
}

type ResourceLock struct {
	logger log.Logger
	conn   *zk.Conn
	clean  func()
}

func NewResourceLock(logger log.Logger, zkServers string) (lock *ResourceLock, err error) {
	servers := strings.Split(zkServers, ",")
	if len(servers) == 0 {
		err = errors.New("zk servers is required")
		return
	}
	conn, event, err := zk.Connect(servers, time.Second)
	if err != nil {
		return nil, err
	}
	// 等待连接成功
	for {
		isConnected := false
		select {
		case connEvent := <-event:
			if connEvent.State == zk.StateConnected {
				isConnected = true
				logger.Infof("connect to zookeeper server success!")
			}
		case <-time.After(time.Second * 3):
			// 3秒仍未连接成功则返回连接超时
			return nil, errors.New("connect to zookeeper server timeout")
		}
		if isConnected {
			break
		}
	}
	lock = &ResourceLock{
		conn: conn,
		clean: func() {
			conn.Close()
		},
	}
	return
}

func (r *ResourceLock) Close() {
	if r.clean != nil {
		r.clean()
	}
}

const (
	FlagPermanent = 0                // 0: 永久保存
	FlagEphemeral = zk.FlagEphemeral // 1: 短暂,session断开则该节点也被删除
)

func (r *ResourceLock) Create(path string, data []byte, flags int32) (string, error) {
	return r.conn.Create(path, data, flags, zk.WorldACL(zk.PermAll))
}

func (r *ResourceLock) Exists(path string) (exists bool, err error) {
	exists, _, err = r.conn.Exists(path)
	return
}

func (r *ResourceLock) Get(path string) (data []byte, err error) {
	data, _, err = r.conn.Get(path)
	return
}

func (r *ResourceLock) Watch(path string) (childCh <-chan zk.Event, err error) {
	_, _, childCh, err = r.conn.ChildrenW(path)
	return
}

func (r *ResourceLock) IsConnected() bool {
	if r.conn == nil || r.conn.State() != zk.StateConnected {
		return false
	}
	return true
}
