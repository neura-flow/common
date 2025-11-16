package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	logger, err := NewLogger(&Config{
		Caller: CallerConfig{
			Enabled: true,
			Skip:    0,
		},
		MessageKey: "mymessage",
		Std: StdConfig{
			Enabled: true,
		},
	})
	assert.NoError(t, err)
	conf := logger.Config()
	assert.Equal(t, conf.TimestampKey, DefaultConfig.TimestampKey)
	assert.Equal(t, conf.MessageKey, "mymessage")
	assert.True(t, conf.Caller.Enabled)
	logger.Infof("test %v %v", "abcd", 111)
	logger = logger.WithOptions()
	logger.Infof("aaa")
}

func TestLevel_LowerThan(t *testing.T) {
	l := LevelInfo
	l1 := LevelDebug
	l2 := LevelError
	assert.True(t, l.LowerThan(l2))
	assert.False(t, l.LowerThan(l1))
}
