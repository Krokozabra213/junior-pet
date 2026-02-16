// Package logger provides environment-based structured logging with dynamic level control.
package logger

import (
	"log/slog"
	"os"
)

const (
	EnvLocal = "local"
	EnvDev   = "development"
	EnvProd  = "prod"
)

type Logger struct {
	*slog.Logger
	level *slog.LevelVar
}

func Init(env string) *Logger {
	l := Setup(env)
	slog.SetDefault(l.Logger)
	return l
}

func Setup(env string) *Logger {
	lvl := &slog.LevelVar{}

	var handler slog.Handler

	switch env {
	case EnvLocal:
		lvl.Set(slog.LevelDebug)
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})

	case EnvDev:
		lvl.Set(slog.LevelDebug)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})

	case EnvProd:
		lvl.Set(slog.LevelInfo)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: false,
		})

	default:
		lvl.Set(slog.LevelInfo)
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})
	}

	return &Logger{
		Logger: slog.New(handler),
		level:  lvl,
	}
}

func (l *Logger) SetLevel(level slog.Level) {
	l.level.Set(level)
	l.Info("log level changed", slog.String("new_level", level.String()))
}

func (l *Logger) GetLevel() slog.Level {
	return l.level.Level()
}

func Err(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "")
	}
	return slog.String("error", err.Error())
}

func (l *Logger) WithService(name, version string) *Logger {
	return &Logger{
		Logger: l.Logger.With(
			slog.String("service", name),
			slog.String("version", version),
		),
		level: l.level,
	}
}
