package log

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/neura-flow/common/metadata"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	cfg         Config
	zapCore     zapcore.Core
	logger      *zap.Logger
	md          metadata.Metadata
	fields      []zap.Field
	level       zap.AtomicLevel
	encoderConf zapcore.EncoderConfig
}

func init() {
	l, err := NewZapLogger(DefaultConfig)
	if err != nil {
		panic(err)
	}
	l.Infof("init default logger OK")
}

func NewZapLogger(cfg *Config) (Logger, error) {
	l := &ZapLogger{
		cfg: *cfg,
		md:  metadata.New(),
	}
	if cfg.Fields != "" {
		for _, f := range strings.Split(cfg.Fields, ".") {
			kv := strings.Split(f, "=")
			if len(kv) == 2 {
				if kv[0] != "" {
					l.md.Set(kv[0], kv[1])
				}
			}
		}
	}
	if l.cfg.Encoding == "" {
		l.cfg.Encoding = DefaultConfig.Encoding
	}
	if l.cfg.MessageKey == "" {
		l.cfg.MessageKey = DefaultConfig.MessageKey
	}
	if l.cfg.TimestampKey == "" {
		l.cfg.TimestampKey = DefaultConfig.TimestampKey
	}

	l.cfg.Level = Level(strings.ToLower(string(l.cfg.Level)))
	l.encoderConf = zap.NewProductionConfig().EncoderConfig
	l.encoderConf.EncodeTime = zapcore.ISO8601TimeEncoder
	l.encoderConf.MessageKey = l.cfg.MessageKey
	l.encoderConf.TimeKey = l.cfg.TimestampKey
	l.level = zap.NewAtomicLevel()
	if level, ok := zapLevels[l.cfg.Level]; ok {
		l.level.SetLevel(level)
	} else {
		l.level.SetLevel(zapcore.InfoLevel)
	}
	l.initCore()
	l.initLogger()
	return l, nil
}

func (l *ZapLogger) Config() Config {
	return l.cfg
}

func (l *ZapLogger) initCore() {
	var writeSyncers []zapcore.WriteSyncer
	if l.cfg.Std.Enabled {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stderr))
	}
	if l.cfg.File.Enabled {
		fileWriter := &lumberjack.Logger{
			Filename: l.cfg.File.Path,
			MaxSize:  l.cfg.File.MaxSize,
			MaxAge:   l.cfg.File.MaxDays,
			Compress: l.cfg.File.Compress,
		}
		writeSyncers = append(writeSyncers, zapcore.AddSync(fileWriter))
	}
	var encoder zapcore.Encoder
	if l.cfg.Encoding == EncodingJSON {
		encoder = zapcore.NewJSONEncoder(l.encoderConf)
	} else {
		encoder = zapcore.NewConsoleEncoder(l.encoderConf)
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(writeSyncers...)
	l.zapCore = zapcore.NewCore(encoder, writeSyncer, l.level)
}

func (l *ZapLogger) initLogger() {
	var defaultSkipLevel = 1 + l.cfg.Caller.Skip
	zapOpts := []zap.Option{
		zap.AddCallerSkip(defaultSkipLevel),
		zap.ErrorOutput(os.Stdout),
		zap.WithCaller(l.cfg.Caller.Enabled),
	}
	if l.cfg.Stack.Enabled {
		zapOpts = append(zapOpts, zap.AddStacktrace(zap.NewAtomicLevelAt(zap.FatalLevel)))
	}
	fs := make([]zap.Field, 0)
	l.md.Range(func(kv metadata.KV) {
		fs = append(fs, zap.Any(kv.Key(), kv.Value()))
	})
	if len(fs) > 0 {
		sort.Slice(fs, func(i, j int) bool {
			return fs[i].Key < fs[j].Key
		})
		zapOpts = append(zapOpts, zap.Fields(fs...))
	}
	zlog := zap.New(l.zapCore, zapOpts...)
	l.logger = zlog
	l.fields = fs
}

func (l *ZapLogger) WithOptions(opts ...Option) Logger {
	zlog := *l
	zlog.md = metadata.Clone(l.md)
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	zlog.doWithOptions(options)
	return &zlog
}

func (l *ZapLogger) With(keyValues ...interface{}) Logger {
	n := len(keyValues) - 1
	kvs := make([]metadata.KV, 0, n)
	for i := 0; i < n; i += 2 {
		key, ok := keyValues[i].(string)
		if !ok {
			continue
		}
		value := keyValues[i+1]
		kvs = append(kvs, metadata.NewKV(key, value))
	}
	zlog := *l
	zlog.WithOptions(WithFields(kvs...))
	return &zlog
}

func (l *ZapLogger) doWithOptions(opts *Options) {
	if len(opts.Fields) > 0 {
		for _, kv := range opts.Fields {
			l.md.Set(kv.Key(), kv.Value())
		}
		l.initLogger()
	}
	if opts.Skip > 0 {
		l.logger = l.logger.WithOptions(zap.AddCallerSkip(opts.Skip))
	}
}

var zapLevels = map[Level]zapcore.Level{
	LevelDebug: zap.DebugLevel,
	LevelInfo:  zap.InfoLevel,
	LevelWarn:  zap.WarnLevel,
	LevelError: zap.ErrorLevel,
	LevelFatal: zap.FatalLevel,
	LevelPanic: zap.PanicLevel,
}

func (l *ZapLogger) IsLevel(level Level) bool {
	return l.cfg.Level == level
}

func (l *ZapLogger) Log(level Level, kvs ...interface{}) error {
	zapLevel, ok := zapLevels[level]
	if !ok {
		zapLevel = zap.InfoLevel
	}
	if zapLevels[l.cfg.Level] > zapLevel {
		return nil
	}
	if len(kvs) == 0 || len(kvs)%2 != 0 {
		l.Warnf("kvs must appear in pairs: %v", kvs)
		return nil
	}

	var fields []zap.Field
	copy(fields, l.fields)
	msg := ""
	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			l.Warnf("key must be string, key: %v", kvs[i])
			continue
		}
		if key == "msg" || key == "message" {
			key = "message"
			msg = fmt.Sprint(kvs[i+1])
		} else {
			fields = append(fields, zap.Any(key, kvs[i+1]))
		}
	}
	if ce := l.logger.Check(zapLevel, msg); ce != nil {
		ce.Write(fields...)
	}
	return nil
}

func (l *ZapLogger) Debugf(msg string, v ...interface{}) {
	l.logger.Debug(fmt.Sprintf(msg, v...))
}

func (l *ZapLogger) Infof(msg string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, v...))
}

func (l *ZapLogger) Warnf(msg string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, v...))
}

func (l *ZapLogger) Errorf(msg string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, v...))
}

func (l *ZapLogger) Fatalf(msg string, v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(msg, v...))
}

func (l *ZapLogger) Panicf(msg string, v ...interface{}) {
	l.logger.Panic(fmt.Sprintf(msg, v...))
}
