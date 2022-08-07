package log

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 日志配置
type Config struct {
	Level       string      `yaml:"level" json:"level"`               // 日志等级, debug,info,warn,error,dpanic,panic,fatal 默认warn
	Format      string      `yaml:"format" json:"format"`             // 编码格式: json or console 默认json
	EncodeLevel string      `yaml:"encode_level" json:"encode_level"` // 编码器类型, 默认 LowercaseLevelEncoder
	Adapter     string      `yaml:"adapter" json:"adapter"`           // 输出: file,console,multi,custom 默认 console
	Stack       bool        `yaml:"stack" json:"stack"`               // 使能栈调试输出 , 默认false
	Path        string      `yaml:"path" json:"path"`                 // 日志存放路径, 默认 empty
	Writer      []io.Writer `yaml:"-" json:"-"`                       // 当 adapter=custom使用,如为writer为空,将使用os.Stdout
	// see lumberjack.Logger
	Filename   string `yaml:"filename" json:"filename"`       // 文件名,空字符使用默认    默认<processname>-lumberjack.log
	MaxSize    int    `yaml:"max_size" json:"max_size"`       // 每个日志文件最大尺寸(MB) 默认100MB,
	MaxAge     int    `yaml:"max_age" json:"max_age"`         // 日志文件保存天数, 默认0不删除
	MaxBackups int    `yaml:"max_backups" json:"max_backups"` // 日志文件保存备份数, 默认0都保存
	LocalTime  bool   `yaml:"local_time" json:"local_time"`   // 是否格式化时间戳, 默认UTC时间
	Compress   bool   `yaml:"compress" json:"compress"`       // 压缩文件,采用gzip, 默认不压缩
}

func New(c Config) *zap.Logger {
	var options []zap.Option
	var encoder zapcore.Encoder

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    toEncodeLevel(c.EncodeLevel),
		EncodeTime:     zapcore.ISO8601TimeEncoder, // 修改输出时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if c.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 设置日志输出等级
	level := zap.NewAtomicLevelAt(toLevel(c.Level))
	// 初始化core
	core := zapcore.NewCore(encoder, toWriter(&c), level)

	// 添加显示文件名和行号,跳过封装调用层,栈调用,及使能等级
	if c.Stack {
		stackLevel := zap.NewAtomicLevel()
		stackLevel.SetLevel(zap.WarnLevel) // 只显示栈的错误等级
		options = append(options,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(stackLevel),
		)
	}
	return zap.New(core, options...)
}

func toLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.WarnLevel
	}
}

func toEncodeLevel(l string) zapcore.LevelEncoder {
	switch l {
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		fallthrough
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func toWriter(c *Config) zapcore.WriteSyncer {
	switch strings.ToLower(c.Adapter) {
	case "file":
		return zapcore.AddSync(&lumberjack.Logger{ // 文件切割
			Filename:   filepath.Join(c.Path, c.Filename),
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			LocalTime:  c.LocalTime,
			Compress:   c.Compress,
		})
	case "multi":
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&lumberjack.Logger{ // 文件切割
			Filename:   filepath.Join(c.Path, c.Filename),
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			LocalTime:  c.LocalTime,
			Compress:   c.Compress,
		}))
	case "custom":
		ws := make([]zapcore.WriteSyncer, 0, len(c.Writer))

		for _, writer := range c.Writer {
			ws = append(ws, zapcore.AddSync(writer))
		}
		if len(ws) == 0 {
			return zapcore.AddSync(os.Stdout)
		}
		if len(ws) == 1 {
			return ws[0]
		}
		return zapcore.NewMultiWriteSyncer(ws...)
	case "console":
		fallthrough
	default:
		return zapcore.AddSync(os.Stdout)
	}
}
