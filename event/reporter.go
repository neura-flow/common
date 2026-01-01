package event

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/metadata"
	"github.com/prometheus/client_golang/prometheus"
)

// Reporter 报告器
type Reporter interface {
	//Report 报告事件的处理耗时和错误
	Report(e Event, latency time.Duration, err Error)
}

// ReportFunc 报告函数
type ReportFunc func(e Event, latency time.Duration, err Error)

// Report 报告事件的处理耗时和错误
func (f ReportFunc) Report(e Event, latency time.Duration, err Error) {
	f(e, latency, err)
}

// Alarm 告警器
type Alarm interface {
	// SetLevel 设置告警器告警等级，低于该等级不告警
	SetLevel(level log.Level)
	// Alarm 发送告警
	Alarm(level log.Level, component string, errCode string, err string, traceId string)
}

// AlarmReporter 创建一个告警报告器，用于根据告警等级决定是否告警
func AlarmReporter(alarm Alarm, level func(err Error) log.Level) Reporter {
	return ReportFunc(func(e Event, latency time.Duration, err Error) {
		if err == nil {
			return
		}
		l := level(err)
		if l != log.LevelFatal && l != log.LevelError && l != log.LevelWarn {
			return
		}
		var traceId = ""
		var component = ""
		if md := e.Metadata(); md != nil {
			traceId = fmt.Sprintf("%v", md.Value(metadata.KeyTraceId))
			component = fmt.Sprintf("%v", md.Value(metadata.KeyComponent))
		}
		alarm.Alarm(level(err), component, string(err.Code()), err.Error(), fmt.Sprintf("%v", traceId))
	})
}

// LogReporter 创建一个日志报告器
func LogReporter(logger log.Logger, checker func(e Event) bool) Reporter {
	return ReportFunc(func(e Event, latency time.Duration, err Error) {
		if !checker(e) {
			return
		}
		logger = logger.WithOptions(log.WithSkip(2))
		if err != nil {
			logger.Errorf("event:%v, latency:%d, error:%v", e.Metadata(), latency, err)
		} else {
			logger.Infof("event:%v, latency:%d", e, latency)
		}
	})
}

// MetricsReporterOpts 监控指标报告器选项
type MetricsReporterOpts struct {
	Namespace   string
	Subsystem   string
	LabelKeys   []string
	ConstLabels map[string]string
}

// MetricsReporter 创建一个监控指标报告器
func MetricsReporter(opts MetricsReporterOpts) (Reporter, error) {
	if opts.Namespace == "" || opts.Subsystem == "" {
		return nil, errors.New("namespace and subsystem can't be empty")
	}

	durationOpts := prometheus.HistogramOpts{
		Namespace:   opts.Namespace,
		Subsystem:   opts.Subsystem,
		Name:        "duration",
		Help:        fmt.Sprintf("Duration(ms) of event in %s_%s", opts.Namespace, opts.Subsystem),
		Buckets:     []float64{10, 25, 50, 100, 250, 500, 1000},
		ConstLabels: opts.ConstLabels,
	}

	counterOpts := prometheus.CounterOpts{
		Namespace:   opts.Namespace,
		Subsystem:   opts.Subsystem,
		Name:        "result",
		Help:        fmt.Sprintf("The total number of processed event in %s_%s", opts.Namespace, opts.Subsystem),
		ConstLabels: opts.ConstLabels,
	}
	labelKeys := append([]string{}, opts.LabelKeys...)
	histogram := prometheus.NewHistogramVec(durationOpts, labelKeys)
	labelKeys = append([]string{}, opts.LabelKeys...)
	labelKeys = append(labelKeys, "code")
	counter := prometheus.NewCounterVec(counterOpts, labelKeys)

	if err := prometheus.Register(counter); err != nil {
		return nil, err
	}
	if err := prometheus.Register(histogram); err != nil {
		return nil, err
	}
	//用内存池避免每次都创建一个用于label的map
	pool := sync.Pool{
		New: func() interface{} {
			return make(map[string]string)
		},
	}
	return ReportFunc(func(e Event, latency time.Duration, err Error) {
		md := e.Metadata()
		if md == nil {
			return
		}
		var labels = pool.Get().(map[string]string)
		defer pool.Put(labels)
		for _, key := range opts.LabelKeys {
			v := md.Value(key)
			if v == nil {
				labels[key] = ""
			} else {
				labels[key] = fmt.Sprintf("%v", v)
			}
		}
		histogram.With(labels).Observe(float64(latency.Milliseconds()))
		if err != nil {
			labels["code"] = string(err.Code())
		} else {
			labels["code"] = ""
		}
		counter.With(labels).Add(1)
	}), nil
}
