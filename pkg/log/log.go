package log

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

// Logger обертка над slog.Logger
type Logger struct {
	logger *slog.Logger
}

var (
	globalLogger *Logger
	once         sync.Once
)

func NewLogger() *Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return &Logger{logger: slog.New(handler)}
}

func getGlobalLogger() *Logger {
	if globalLogger == nil {
		once.Do(func() {
			globalLogger = NewLogger()
		})
	}
	return globalLogger
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fieldsToAny(fields)...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fieldsToAny(fields)...)
}

func (l *Logger) Error(err string, fields ...Field) {
	l.logger.Error(err, fieldsToAny(fields)...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	all := append(extractFields(ctx), fieldsToAny(fields)...)
	l.logger.Info(msg, all...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, fields ...Field) {
	all := append(extractFields(ctx), fieldsToAny(fields)...)
	l.logger.Debug(msg, all...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	all := append(extractFields(ctx), fieldsToAny(fields)...)
	l.logger.Error(msg, all...)
}

func Info(msg string, fields ...Field) {
	getGlobalLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	getGlobalLogger().Warn(msg, fields...)
}

func Error(err string, fields ...Field) {
	getGlobalLogger().Error(err, fields...)
}

func InfoContext(ctx context.Context, msg string, fields ...Field) {
	getGlobalLogger().InfoContext(ctx, msg, fields...)
}

func DebugContext(ctx context.Context, msg string, fields ...Field) {
	getGlobalLogger().DebugContext(ctx, msg, fields...)
}

func ErrorContext(ctx context.Context, msg string, fields ...Field) {
	getGlobalLogger().ErrorContext(ctx, msg, fields...)
}
