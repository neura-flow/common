package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/types"
)

type MetricsConfig struct {
	Enabled        bool   `json:"enabled,omitempty"`        // 是否开启监控
	Keys           string `json:"keys,omitempty"`           // 需要监控的一些key(模糊匹配)
	SlowLogMinCost int    `json:"slowLogMinCost,omitempty"` // 慢日志最低耗时, <=0表示关闭

	clusterId string `json:"-"` // 集群ID, 取shortName
}

type Config struct {
	//redis地址，多个地址用逗号隔开
	Addrs    string        `json:"addrs,omitempty"`
	Password string        `json:"password,omitempty"`
	DB       int           `json:"db,omitempty"`
	Timeout  types.Timeout `json:"timeout,omitempty"`
	Pool     types.Pool    `json:"pool,omitempty"`
	Ping     bool          `json:"ping,omitempty"`
	// Kind redis类型,支持simple、cluster、failover
	Kind    string        `json:"kind,omitempty"`
	Metrics MetricsConfig `json:"metrics,omitempty"`
}

type Client struct {
	redis.UniversalClient
	cfg *Config
}

func NewClient(ctx context.Context, logger log.Logger, cfg *Config) (*Client, error) {
	var client redis.UniversalClient

	options := &redis.UniversalOptions{
		Addrs:         strings.Split(cfg.Addrs, ","),
		Password:      cfg.Password,
		DB:            cfg.DB,
		DialTimeout:   time.Duration(cfg.Timeout.Dail) * time.Millisecond,
		ReadTimeout:   time.Duration(cfg.Timeout.Read) * time.Millisecond,
		WriteTimeout:  time.Duration(cfg.Timeout.Write) * time.Millisecond,
		PoolSize:      cfg.Pool.MaxOpen,
		MinIdleConns:  cfg.Pool.MinIdle,
		MaxConnAge:    time.Duration(cfg.Pool.LifeTime) * time.Second,
		PoolFIFO:      true,
		RouteRandomly: true,
	}

	switch cfg.Kind {
	case "simple":
		client = redis.NewClient(options.Simple())
	case "cluster":
		client = redis.NewClusterClient(options.Cluster())
	case "failover":
		client = redis.NewFailoverClient(options.Failover())
	default:
		return nil, fmt.Errorf("invalid redis kind: %s", cfg.Kind)
	}

	if cfg.Ping {
		cmd := client.Ping(ctx)
		if err := cmd.Err(); err != nil {
			client.Close()
			return nil, err
		}
	}
	if cfg.Metrics.Enabled {
		cfg.Metrics.clusterId = ""
		client.AddHook(NewHook(cfg, logger))
	}
	cli := &Client{
		UniversalClient: client,
		cfg:             cfg,
	}
	return cli, nil
}
