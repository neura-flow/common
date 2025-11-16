package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/types"
	"gorm.io/driver/mysql"
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
	//Addr 地址
	Addr string `json:"addr,omitempty"`
	//Username 帐号
	Username string `json:"username,omitempty"`
	//Password 密码
	Password string `json:"password,omitempty"`
	//DB 数据库
	DB string `json:"db,omitempty"`
	//Timeout 超时参数
	Timeout types.Timeout `json:"timeout,omitempty"`
	//Pool 连接池参数
	Pool types.Pool `json:"pool,omitempty"`
	//Check 是否开启健康检查
	Check bool `json:"check,omitempty"`
	//Options 其他参数，格式 aaa=bbb&ccc=ddd，具体参数参考 gorm 文档
	Options string `json:"options,omitempty"`
	// mysql监控
	Metrics MetricsConfig `json:"metrics,omitempty"`
}

func (cfg *Config) DSN() string {
	check := "False"
	if cfg.Check {
		check = "True"
	}
	params := fmt.Sprintf("timeout=%dms&readTimeout=%dms&writeTimeout=%dms&checkConnLiveness=%s",
		cfg.Timeout.Dail, cfg.Timeout.Read, cfg.Timeout.Write, check)
	if cfg.Options != "" {
		params = params + "&" + cfg.Options
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.Username, cfg.Password, cfg.Addr, cfg.DB, params)
	return dsn
}

type Client struct {
	*gorm.DB
	cfg *Config
}

func NewClient(ctx context.Context, logger log.Logger, cfg *Config) (*Client, error) {
	gormConfig := &gorm.Config{
		Logger: &loggerAdapter{logger},
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
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
