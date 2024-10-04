package log

import (
	"context"
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

func init() {
	singletonMap.AddFactory(defaultKey, defaultInstance)
}

func AddFactory(key string, factory func() *Logger) {
	singletonMap.AddFactory(key, factory)
}

func Remove(key string) {
	singletonMap.Remove(key)
}

func Get(key string) *Logger {
	return singletonMap.GetInstance(key)
}

func defaultInstance() *Logger {
	options := Options{
		Filename:   "log/log.log",
		MaxBackup:  365,
		IsCompress: false,
	}
	return NewLogger(options)
}

type Options struct {
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

	slogger := slog.New(slog.NewJSONHandler(loggingFile, &slog.HandlerOptions{
		AddSource: true,
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
	}))

	return &Logger{
		slogger: slogger.With("PID", os.Getpid()),
	}
}

func Default() *Logger {
	return singletonMap.GetInstance(defaultKey)
}

type Logger struct {
	slogger *slog.Logger
}

func (l Logger) InfoJson(msg string, attrs ...slog.Attr) {
	l.slogger.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

func (l Logger) ErrorJson(msg string, attrs ...slog.Attr) {
	l.slogger.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}

func (l Logger) Info(msg string, args ...any) {
	l.slogger.Info(msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	l.slogger.Error(msg, args...)
}
