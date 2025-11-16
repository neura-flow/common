package sqlite

import (
	"context"
	"time"

	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MetricsConfig struct {
	Enabled        bool   `json:"enabled,omitempty"`        // 是否开启监控
	Tables         string `json:"tables,omitempty"`         // 要监控的表,用逗号隔开
	SlowLogMinCost int    `json:"slowLogMinCost,omitempty"` // 慢日志最低耗时, <=0表示关闭
	clusterId      string `json:"-"`                        // 集群ID
}

type Config struct {
	File string `json:"file"`
	//Timeout 超时参数
	Timeout types.Timeout `json:"timeout,omitempty"`
	//Pool 连接池参数
	Pool types.Pool `json:"pool,omitempty"`
	// mysql监控
	Metrics MetricsConfig `json:"metrics,omitempty"`
}

type Client struct {
	*gorm.DB
	cfg *Config
}

func NewClient(ctx context.Context, logger log.Logger, cfg *Config) (*Client, error) {
	gormConfig := &gorm.Config{
		Logger: &loggerAdapter{logger},
	}
	db, err := gorm.Open(sqlite.Open(cfg.File), gormConfig)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Pool.LifeTime) * time.Second)
	sqlDB.SetMaxIdleConns(cfg.Pool.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.Pool.MaxOpen)

	if cfg.Metrics.Enabled {
		cfg.Metrics.clusterId = ""
		initializeMetrics(db, cfg, logger)
	}

	cli := &Client{
		DB:  db,
		cfg: cfg,
	}
	return cli, nil
}

type loggerAdapter struct {
	logger log.Logger
}

func (a *loggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	return a
}
func (a *loggerAdapter) Info(ctx context.Context, message string, params ...interface{}) {
	a.logger.Infof("message: %s params: %v", message, params)
}
func (a *loggerAdapter) Warn(ctx context.Context, message string, params ...interface{}) {
	a.logger.Infof("message: %s params: %v", message, params)
}
func (a *loggerAdapter) Error(ctx context.Context, message string, params ...interface{}) {
	a.logger.Infof("message: %s params: %v", message, params)
}
func (a *loggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

}
