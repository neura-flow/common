package redis

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/neura-flow/common/log"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestHook(t *testing.T) {
	var cfg = &Config{
		Metrics: MetricsConfig{
			Enabled:   true,
			Keys:      "foo, aha",
			clusterId: "test_redis",
		},
	}

	assert := assert.New(t)
	logger := log.DefaultLogger()

	t.Run("create a new hook", func(t *testing.T) {
		hook := NewHook(cfg, logger)
		assert.NotNil(hook)
	})

	t.Run("do not panic if metrics are already registered", func(t *testing.T) {
		NewHook(cfg, logger)
		assert.NotPanics(func() {
			NewHook(cfg, logger)
		})
	})

	t.Run("export metrics after a command is processed", func(t *testing.T) {
		hook := NewHook(cfg, logger)

		cmd := redis.NewStringCmd(context.TODO(), "get", "foo")
		cmd.SetErr(errors.New("some error"))

		ctx, err1 := hook.BeforeProcess(context.TODO(), cmd)
		err2 := hook.AfterProcess(ctx, cmd)

		assert.Nil(err1)
		assert.Nil(err2)

		metrics, err := prometheus.DefaultGatherer.Gather()
		assert.Nil(err)

		assert.ElementsMatch([]string{
			"redis_client_duration",
			"redis_client_result",
			"redis_client_size",
		}, filter(metrics, "redis_client"))
	})

	t.Run("export metrics after a pipeline is processed", func(t *testing.T) {
		hook := NewHook(cfg, logger)

		cmd1 := redis.NewStringCmd(context.TODO(), "get foo")
		cmd1.SetErr(errors.New("some error"))
		cmd2 := redis.NewStringCmd(context.TODO(), "get haha")
		cmd2.SetErr(nil)

		cmds := []redis.Cmder{cmd1, cmd2}
		ctx, err1 := hook.BeforeProcessPipeline(context.TODO(), cmds)
		err2 := hook.AfterProcessPipeline(ctx, cmds)

		assert.Nil(err1)
		assert.Nil(err2)

		metrics, err := prometheus.DefaultGatherer.Gather()
		assert.Nil(err)

		assert.ElementsMatch([]string{
			"redis_client_result",
			"redis_client_duration",
			"redis_client_size",
		}, filter(metrics, "redis_client"))
	})
}

func filter(metrics []*io_prometheus_client.MetricFamily, namespace string) []string {
	var result = make([]string, 0)
	for _, metric := range metrics {
		if strings.HasPrefix(*metric.Name, namespace) {
			result = append(result, *metric.Name)
		}
	}
	return result
}
