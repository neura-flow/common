package log

const (
	EncodingJSON    = "json"
	EncodingConsole = "console"

	LevelDebug = Level("debug")
	LevelInfo  = Level("info")
	LevelWarn  = Level("warn")
	LevelError = Level("error")
	LevelFatal = Level("fatal")
	LevelPanic = Level("panic")

	fieldKeyError = "error"
)

type Level string

var levels = map[Level]int{
	LevelDebug: 0,
	LevelInfo:  1,
	LevelWarn:  2,
	LevelError: 3,
	LevelFatal: 4,
	LevelPanic: 5,
}

func (l Level) LowerThan(l1 Level) bool {
	return levels[l] < levels[l1]
}

var DefaultConfig = &Config{
	Std: StdConfig{
		Enabled: true,
	},
	Level:    LevelInfo,
	Encoding: EncodingJSON,
	Caller: CallerConfig{
		Enabled: true,
	},
	MessageKey:   "message",
	TimestampKey: "timestamp",
}

type CallerConfig struct {
	Enabled bool `json:"enabled,omitempty" desc:"是否开启调用者，有性能损耗，默认不开启"`
	Skip    int  `json:"skip,omitempty" desc:"跳过几层调用栈"`
}

type StackConfig struct {
	Enabled bool `json:"enabled,omitempty" desc:"是否开启调用栈，有性能损耗，默认不开启"`
}

type FileConfig struct {
	Enabled  bool   `json:"enabled,omitempty" desc:"是否开启日志文件，默认不开启"`
	Path     string `json:"path,omitempty" desc:"日志文件路径"`
	Compress bool   `json:"compress,omitempty" desc:"是否开启压缩，默认不开启"`
	MaxSize  int    `json:"maxSize,omitempty" desc:"文件大小最大值"`
	MaxDays  int    `json:"maxDays,omitempty" desc:"文件最大保存天数"`
}

type StdConfig struct {
	Enabled bool `json:"enabled,omitempty" desc:"是否开启标准输出"`
}

type Config struct {
	Level        Level        `json:"level,omitempty" desc:"日志级别，默认 info 级别"`
	Encoding     string       `json:"encoding,omitempty" desc:"日志格式，默认 json"`
	File         FileConfig   `json:"file,omitempty" desc:"日志文件配置项"`
	Std          StdConfig    `json:"std,omitempty" desc:"标准输出配置项"`
	Caller       CallerConfig `json:"caller,omitempty" desc:"调用者配置项"`
	Stack        StackConfig  `json:"stack,omitempty" desc:"调用栈配置项"`
	MessageKey   string       `json:"messageKey,omitempty" desc:"message字段key名称，默认 message"`
	TimestampKey string       `json:"timestampKey,omitempty" desc:"timestamp字段key名称，默认timestamp"`
	Fields       string       `json:"fields,omitempty" desc:"日志字段，格式：key1=val1,key2=val2"`
}
