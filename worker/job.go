package worker

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/neura-flow/common/event"
	"github.com/neura-flow/common/exception"
	"github.com/neura-flow/common/metadata"
)

type Job interface {
	event.Event
	Do()
}

const (
	JobTypeKey = "jobType"

	JobTypeFunc = iota
	JobTypeContext
	JobTypeDelay
	JobTypeCron
)

type FuncJob func()

func (f FuncJob) Do() {
	f()
}

func (f FuncJob) Metadata() metadata.Metadata {
	return metadata.FromMap(map[string]interface{}{JobTypeKey: JobTypeFunc})
}

func (f FuncJob) Data() interface{} {
	return nil
}

type ContextJob struct {
	Ctx context.Context
	F   func(ctx context.Context)
}

func (j ContextJob) Do() {
	j.F(j.Ctx)
}

func (j ContextJob) Metadata() metadata.Metadata {
	md := metadata.FromMap(map[string]interface{}{JobTypeKey: JobTypeContext})
	return metadata.MergeMetadata(metadata.FromContext(j.Ctx), md)
}

func (j ContextJob) Data() interface{} {
	return j.Ctx
}

type DelayJob struct {
	Job
	Delay time.Duration
}

func (j DelayJob) Metadata() metadata.Metadata {
	md := metadata.FromMap(map[string]interface{}{JobTypeKey: JobTypeDelay})
	return metadata.MergeMetadata(j.Job.Metadata(), md)
}

type CronJob struct {
	Job
	Interval time.Duration
}

func (j CronJob) Metadata() metadata.Metadata {
	md := metadata.FromMap(map[string]interface{}{JobTypeKey: JobTypeCron})
	return metadata.MergeMetadata(j.Job.Metadata(), md)
}

// JobGroup 任务组
type JobGroup struct {
	cancelOnError bool
	ctx           context.Context
	cancel        func()
	wg            *sync.WaitGroup
	lock          sync.RWMutex
	errs          jobErrors
	recover       bool
}

// JobGroupError 任务组错误
type JobGroupError interface {
	error
	// JobError 获取某个任务的错误
	JobError(key interface{}) error
}

type jobErrors map[interface{}]error

func (je jobErrors) Error() string {
	data, _ := json.Marshal(je)
	return string(data)
}

func (je jobErrors) JobError(key interface{}) error {
	if je == nil {
		return nil
	}
	return je[key]
}

// JobGroupOption 任务组参数选项
type JobGroupOption func(o *JobGroup)

// JobGroupWithContext 使用 context ，可用于超时控制、主动取消、传递元数据等
func JobGroupWithContext(ctx context.Context) JobGroupOption {
	return func(o *JobGroup) {
		o.ctx = ctx
	}
}

// JobGroupWithCancelOnError 错误时取消所有任务
func JobGroupWithCancelOnError() JobGroupOption {
	return func(o *JobGroup) {
		o.cancelOnError = true
	}
}

// JobGroupWithRecover 捕获异常
func JobGroupWithRecover() JobGroupOption {
	return func(o *JobGroup) {
		o.recover = true
	}
}

// NewJobGroup 创建一个任务组，注意：一旦 Wait 或者 Cancel 后不可重复使用
func NewJobGroup(opts ...JobGroupOption) *JobGroup {
	jg := &JobGroup{
		wg:   &sync.WaitGroup{},
		lock: sync.RWMutex{},
		errs: map[interface{}]error{},
	}
	for _, opt := range opts {
		opt(jg)
	}
	if jg.ctx == nil {
		jg.ctx = context.Background()
	}
	jg.ctx, jg.cancel = context.WithCancel(jg.ctx)
	return jg
}

// NewJob 用任务组创建一个任务
func (g *JobGroup) NewJob(f func() error, key interface{}) Job {
	g.wg.Add(1)
	return FuncJob(func() {
		var err error
		defer func() {
			g.wg.Done()
			if err != nil {
				if g.cancelOnError {
					g.Cancel()
				}
				g.lock.Lock()
				g.errs[key] = err
				g.lock.Unlock()
			}
		}()
		select {
		case <-g.ctx.Done():
			return
		default:
		}
		defer exception.Recover(func(pe exception.PanicException) bool {
			err = pe
			return g.recover
		})
		err = f()
	})
}

// Cancel 取消任务组
func (g *JobGroup) Cancel() {
	if g.cancel != nil {
		g.cancel()
	}
}

// CancelAndWait 取消任务组，并等待所有任务取消
func (g *JobGroup) CancelAndWait() JobGroupError {
	g.Cancel()
	return g.Wait()
}

// Wait 等待任务组
func (g *JobGroup) Wait() JobGroupError {
	g.wg.Wait()
	g.lock.Lock()
	defer g.lock.Unlock()
	errs := g.errs
	g.errs = map[interface{}]error{}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
