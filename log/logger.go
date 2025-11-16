package log

import (
	"context"

	"github.com/neura-flow/common/metadata"
)

type Options struct {
	Skip   int
	Fields []metadata.KV
}

type Option func(o *Options)

func WithSkip(skip int) Option {
	return func(o *Options) {
		o.Skip = skip
	}
}

func WithFields(fields ...metadata.KV) Option {
	return func(o *Options) {
		o.Fields = fields
	}
}

type Logger interface {
	Config() Config
	WithOptions(opts ...Option) Logger
	With(keyValues ...interface{}) Logger
	Log(level Level, kvs ...interface{}) error
	Debugf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	Panicf(format string, v ...interface{})
}

func NewLogger(cfg *Config) (Logger, error) {
	return NewZapLogger(cfg)
}

func DefaultLogger() Logger {
	logger, err := NewLogger(&Config{
		Caller: CallerConfig{
			Enabled: true,
			Skip:    0,
		},
		MessageKey: "",
		Std: StdConfig{
			Enabled: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return logger
}

type contextKey struct{}

func FromContext(ctx context.Context) Logger {
	v := ctx.Value(contextKey{})
	if l, ok := v.(Logger); ok {
		return l
	}
	return DefaultLogger()
}

func ToContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}
