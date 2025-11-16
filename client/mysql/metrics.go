package mysql

import (
	"strings"
	"sync"
	"time"

	"github.com/neura-flow/common/log"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

const startTimeKey = "startTime"

var once = sync.Once{}

func initializeMetrics(db *gorm.DB, cfg *Config, logger log.Logger) {
	once.Do(func() {
		m := &metrics{
			cfg:    cfg,
			logger: logger,
			tables: make(map[string]uint8),
		}
		m.init(db)
	})
}

type metrics struct {
	cfg            *Config
	logger         log.Logger
	tables         map[string]uint8
	metricDuration *prometheus.HistogramVec
	metricCounter  *prometheus.CounterVec
}

func (m *metrics) init(db *gorm.DB) {
	m.tables = make(map[string]uint8)
	for _, item := range m.getTables(m.cfg.Metrics.Tables) {
		m.tables[item] = 1
	}

	// 注册监控指标
	metricDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "mysql",
		Subsystem: "client",
		Name:      "duration",
		Help:      "client requests duration(ms).",
		Buckets:   []float64{10, 50, 100, 500, 2000},
	}, []string{"cluster_id", "action", "table"})

	metricCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "mysql",
		Subsystem: "client",
		Name:      "result",
		Help:      "The result of processed requests",
	}, []string{"cluster_id", "action", "table", "success", "msg"})

	m.metricCounter = m.Register(metricCounter).(*prometheus.CounterVec)
	m.metricDuration = m.Register(metricDuration).(*prometheus.HistogramVec)

	// 注入gorm回调函数
	db.Callback().Query().Before("gorm:query").Register("metrics:before_query", m.Before)
	db.Callback().Query().After("gorm:query").Register("metrics:after_query", func(db *gorm.DB) {
		m.Collect(db, "query")
	})

	db.Callback().Create().Before("gorm:create").Register("metrics:before_create", m.Before)
	db.Callback().Create().After("gorm:create").Register("metrics:after_create", func(db *gorm.DB) {
		m.Collect(db, "create")
	})

	db.Callback().Update().Before("gorm:update").Register("metrics:before_update", m.Before)
	db.Callback().Update().After("gorm:update").Register("metrics:after_update", func(db *gorm.DB) {
		m.Collect(db, "update")
	})

	db.Callback().Delete().Before("gorm:delete").Register("metrics:before_delete", m.Before)
	db.Callback().Delete().After("gorm:delete").Register("metrics:after_delete", func(db *gorm.DB) {
		m.Collect(db, "delete")
	})
}

func (m *metrics) getTables(str string) []string {
	var tables = make([]string, 0)
	var arr = strings.Split(str, ",")
	for _, item := range arr {
		tables = append(tables, strings.TrimSpace(item))
	}
	return tables
}

func (m *metrics) Register(collector prometheus.Collector) prometheus.Collector {
	if err := prometheus.Register(collector); err != nil {
		if arErr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return arErr.ExistingCollector
		} else {
			m.logger.Errorf("unexpected error: %s", err.Error())
		}
	}
	return collector
}

func (m *metrics) Before(db *gorm.DB) {
	db.Set(startTimeKey, time.Now())
}

func (m *metrics) GetStartTime(db *gorm.DB) (time.Time, bool) {
	i, ok := db.Get(startTimeKey)
	if !ok {
		return time.Time{}, false
	}
	t, ok := i.(time.Time)
	return t, ok
}

func (m *metrics) Collect(db *gorm.DB, action string) error {
	startTime, ok := m.GetStartTime(db)
	if !ok {
		return nil
	}
	var d = time.Since(startTime).Milliseconds()
	if m.cfg.Metrics.SlowLogMinCost > 0 && d >= int64(m.cfg.Metrics.SlowLogMinCost) {
		var sql = db.Statement.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
		m.logger.Warnf("MySQLSlowLog Latency: %dms, SQL: %s", d, sql)
	}

	if m.tables[db.Statement.Table] == 0 {
		return nil
	}
	var success = "1"
	var msg = ""
	if db.Error != nil {
		msg = db.Error.Error()
		if db.Error != gorm.ErrRecordNotFound {
			success = "0"
		}
	}
	m.metricCounter.WithLabelValues(m.cfg.Metrics.clusterId, action, db.Statement.Table, success, msg).Inc()
	m.metricDuration.WithLabelValues(m.cfg.Metrics.clusterId, action, db.Statement.Table).Observe(float64(d))
	return nil
}
