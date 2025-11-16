package redis

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/neura-flow/common/log"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	Hook struct {
		cfg               *Config
		keys              []string
		logger            log.Logger
		successCollector  *prometheus.CounterVec   // 统计请求是否成功
		sizeCollector     *prometheus.CounterVec   // 统计缓存传输数据量
		durationCollector *prometheus.HistogramVec // 统计请求处理时间
	}

	startKey struct{}
)

// NewHook creates a new go-redis hook instance and registers Prometheus collectors.
func NewHook(cfg *Config, logger log.Logger) *Hook {
	var keys = make([]string, 0)
	var arr = strings.Split(cfg.Metrics.Keys, ",")
	for _, item := range arr {
		keys = append(keys, strings.TrimSpace(item))
	}

	var durationCollector = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "redis",
		Subsystem: "client",
		Name:      "duration",
		Help:      "redis client duration(ms).",
		Buckets:   []float64{5, 10, 50, 100, 500},
	}, []string{"cluster_id", "command", "key"})

	var successCollector = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "redis",
		Subsystem: "client",
		Name:      "result",
		Help:      "The result of processed requests",
	}, []string{"cluster_id", "command", "key", "success", "msg"})

	var sizeCollector = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "redis",
		Subsystem: "client",
		Name:      "size",
		Help:      "The total response size(byte) of processed requests",
	}, []string{"cluster_id", "command", "key"})

	var hook = &Hook{
		cfg:    cfg,
		keys:   keys,
		logger: logger,
	}
	hook.successCollector = hook.register(successCollector).(*prometheus.CounterVec)
	hook.sizeCollector = hook.register(sizeCollector).(*prometheus.CounterVec)
	hook.durationCollector = hook.register(durationCollector).(*prometheus.HistogramVec)
	return hook
}

func (hook *Hook) register(collector prometheus.Collector) prometheus.Collector {
	if err := prometheus.Register(collector); err != nil {
		if arErr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return arErr.ExistingCollector
		} else {
			hook.logger.Errorf("unexpected error: %s", err.Error())
		}
	}
	return collector
}

func (hook *Hook) match(cmd redis.Cmder) (key string, match bool) {
	var str = cmd.String()
	var n = len(hook.keys)
	for i := 0; i < n; i++ {
		if strings.Contains(str, hook.keys[i]) {
			key = hook.keys[i]
			match = true
			break
		}
	}
	return
}

func (hook *Hook) getSize(cmd redis.Cmder) int {
	switch cmd.(type) {
	case *redis.StringCmd:
		return len(cmd.(*redis.StringCmd).Val())
	case *redis.IntCmd:
		return len(strconv.Itoa(int(cmd.(*redis.IntCmd).Val())))
	case *redis.SliceCmd:
		var num = 0
		for _, item := range cmd.(*redis.SliceCmd).Val() {
			if v, ok := item.(string); ok {
				num += len(v)
			} else if v, ok := item.(int64); ok {
				num += len(strconv.Itoa(int(v)))
			}
		}
		return num
	case *redis.IntSliceCmd:
		var num = 0
		for _, item := range cmd.(*redis.IntSliceCmd).Val() {
			num += len(strconv.Itoa(int(item)))
		}
		return num
	case *redis.StringSliceCmd:
		var num = 0
		for _, item := range cmd.(*redis.StringSliceCmd).Val() {
			num += len(item)
		}
		return num
	case *redis.StringStringMapCmd:
		var num = 0
		for k, v := range cmd.(*redis.StringStringMapCmd).Val() {
			num += len(k) + len(v)
		}
		return num
	default:
		return 0
	}
}

func (hook *Hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (hook *Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	startTime, ok := ctx.Value(startKey{}).(time.Time)
	if !ok {
		return nil
	}

	duration := time.Since(startTime).Milliseconds()
	if hook.cfg.Metrics.SlowLogMinCost > 0 && duration >= int64(hook.cfg.Metrics.SlowLogMinCost) {
		if arr := strings.Split(cmd.String(), ":"); len(arr) > 0 {
			hook.logger.Warnf("RedisSlowLog Latency: %dms, Command: %s", duration, arr[0])
		}
	}

	key, match := hook.match(cmd)
	var msg = ""
	var success = "1"
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			msg = "miss"
		} else {
			success = "0"
			msg = err.Error()
		}
	}
	if !match {
		return nil
	}

	var command = cmd.Name()
	hook.sizeCollector.WithLabelValues(hook.cfg.Metrics.clusterId, command, key).Add(float64(hook.getSize(cmd)))
	hook.successCollector.WithLabelValues(hook.cfg.Metrics.clusterId, command, key, success, msg).Inc()
	hook.durationCollector.WithLabelValues(hook.cfg.Metrics.clusterId, command, key).Observe(float64(duration))
	return nil
}

func (hook *Hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startKey{}, time.Now()), nil
}

func (hook *Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	startTime, ok := ctx.Value(startKey{}).(time.Time)
	if !ok {
		return nil
	}

	var ctx1 = context.WithValue(ctx, startKey{}, startTime)
	for i, _ := range cmds {
		hook.AfterProcess(ctx1, cmds[i])
	}

	return nil
}
