package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type Logger struct {
	internal *slog.Logger
}

var (
	globalLogger *Logger
	once         sync.Once
)

func Must() *Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return &Logger{
		internal: slog.New(handler),
	}
}

func getGlobal() *Logger {
	once.Do(func() {
		globalLogger = Must()
	})
	return globalLogger
}

type ctxKeyType struct{}

var ctxKey ctxKeyType

func (l *Logger) LogCtx(ctx context.Context, fields ...any) context.Context {
	existing, _ := ctx.Value(ctxKey).([]any)
	combined := make([]any, 0, len(existing)+len(fields))
	combined = append(combined, existing...)
	combined = append(combined, fields...)
	return context.WithValue(ctx, ctxKey, combined)
}

func (l *Logger) withContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return l.internal
	}

	fields, _ := ctx.Value(ctxKey).([]any)
	if len(fields) == 0 {
		return l.internal
	}

	return l.internal.With(fields...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string) {
	l.withContext(ctx).Info(msg)
}

func (l *Logger) WarnContext(ctx context.Context, msg string) {
	l.withContext(ctx).Warn(msg)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string) {
	l.withContext(ctx).Error(msg)
}

func (l *Logger) InfofContext(ctx context.Context, format string, args ...any) {
	l.withContext(ctx).Info(fmt.Sprintf(format, args...))
}

func (l *Logger) WarnfContext(ctx context.Context, format string, args ...any) {
	l.withContext(ctx).Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) ErrorfContext(ctx context.Context, format string, args ...any) {
	l.withContext(ctx).Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Info(msg string) {
	l.internal.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.internal.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.internal.Error(msg)
}

func (l *Logger) Infof(format string, args ...any) {
	l.internal.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...any) {
	l.internal.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.internal.Error(fmt.Sprintf(format, args...))
}

func InfoContext(ctx context.Context, msg string) {
	getGlobal().InfoContext(ctx, msg)
}

func WarnContext(ctx context.Context, msg string) {
	getGlobal().WarnContext(ctx, msg)
}

func ErrorContext(ctx context.Context, msg string) {
	getGlobal().ErrorContext(ctx, msg)
}

func InfofContext(ctx context.Context, format string, args ...any) {
	getGlobal().InfofContext(ctx, format, args...)
}

func WarnfContext(ctx context.Context, format string, args ...any) {
	getGlobal().WarnfContext(ctx, format, args...)
}

func ErrorfContext(ctx context.Context, format string, args ...any) {
	getGlobal().ErrorfContext(ctx, format, args...)
}

func Info(msg string) {
	getGlobal().Info(msg)
}

func Warn(msg string) {
	getGlobal().Warn(msg)
}

func Error(msg string) {
	getGlobal().Error(msg)
}

func Infof(format string, args ...any) {
	getGlobal().Infof(format, args...)
}

func Warnf(format string, args ...any) {
	getGlobal().Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	getGlobal().Errorf(format, args...)
}
