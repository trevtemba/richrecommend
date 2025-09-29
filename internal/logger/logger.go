package logger

import (
	"log/slog"
	"os"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type LogType string

const (
	LogTypeAgentStart  LogType = "agent_start"
	LogTypeAgentWork   LogType = "agent_work"
	LogTypeAgentFinish LogType = "agent_finish"
	LogTypeAgentAbort  LogType = "agent_abort"
	LogTypeRequest     LogType = "request"
	LogTypeSystem      LogType = "system"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

func Init(handler slog.Handler) {
	log = slog.New(handler)
	slog.SetDefault(log)
}

func Log(logType LogType, level Level, msg string, args ...any) {
	fields := append([]any{"type", logType}, args...)

	switch level {
	case LevelDebug:
		log.Debug(msg, fields...)
	case LevelInfo:
		log.Info(msg, fields...)
	case LevelWarn:
		log.Warn(msg, fields...)
	case LevelError:
		log.Error(msg, fields...)
	default:
		log.Info(msg, fields...)
	}
}
