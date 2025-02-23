package log

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/tingin/base/patterns/singleton"
)

func formatTime(a slog.Attr) slog.Attr {
	if t, ok := a.Value.Any().(time.Time); ok {
		a.Value = slog.StringValue(t.Format(time.DateTime))
	}
	return a
}

func formatSource(a slog.Attr) slog.Attr {
	if src, ok := a.Value.Any().(*slog.Source); ok {
		shortPath := path.Base(src.File)
		shortPath += fmt.Sprintf(":%d", src.Line)
		a.Value = slog.StringValue(shortPath)
	}
	return a
}

var singletonMap = singleton.NewSingletonMap[string, Logger]()

var defaultKey = "Default"

func AddFactory(key string, factory *Logger) {
	singletonMap.AddFactory(key, factory)
}

func Remove(key string) {
	singletonMap.Remove(key)
}

func Get(key string) *Logger {
	return singletonMap.GetInstance(key)
}

func defaultInstance(o Options) *Logger {
	options := Options{
		Filename:   "log/log.log",
		MaxBackup:  365,
		IsCompress: false,
		Json:       true,
		Source:     false,
	}
	if o.Filename != "" {
		options.Filename = o.Filename
	}
	if o.MaxBackup != 0 {
		options.MaxBackup = o.MaxBackup
	}

	options.IsCompress = o.IsCompress
	options.Json = o.Json
	options.Source = o.Source

	return NewLogger(options)
}

type Options struct {
	Json       bool
	Source     bool
	Filename   string
	MaxBackup  int
	IsCompress bool
}

func NewLogger(options Options) *Logger {
	loggingFile := logger.NewLogger().SetOption(&logger.Option{
		FileOption: &logger.FileTimeMode{
			Filename:   options.Filename,
			Maxbuckup:  options.MaxBackup,
			IsCompress: options.IsCompress,
			Timemode:   logger.MODE_DAY,
		},
	},
	)

	var log *slog.Logger

	option := &slog.HandlerOptions{
		AddSource: options.Source,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return formatTime(a)
			case slog.SourceKey:
				return formatSource(a)
			default:
				return a
			}
		},
	}

	if options.Json {
		log = slog.New(slog.NewJSONHandler(loggingFile, option))
	} else {
		log = slog.New(slog.NewTextHandler(loggingFile, option))
	}

	return &Logger{
		slogger: log.With("PID", os.Getpid()),
	}
}

func Default(opt Options) *Logger {
	o := singletonMap.GetInstance(defaultKey)
	if o == nil {
		o = defaultInstance(opt)
		AddFactory(defaultKey, o)
	}
	return singletonMap.GetInstance(defaultKey)
}

type Logger struct {
	slogger *slog.Logger
}

func (l Logger) Info(msg string, args ...any) {
	l.slogger.Info(msg, args...)
}

func (l Logger) Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.slogger.Info(msg)
}

func (l Logger) Error(msg string, args ...any) {
	l.slogger.Error(msg, args...)
}

func (l Logger) Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.slogger.Error(msg)
}
