// nolint
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

func MustNewLogger() *Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return &Logger{logger: slog.New(handler)}
}

func getGlobalLogger() *Logger {
	if globalLogger == nil {
		once.Do(func() {
			globalLogger = MustNewLogger()
		})
	}
	return globalLogger
}

func (l *Logger) InfoContext(ctx context.Context, msg string) {
	l.logger.Info(msg, extractFields(ctx))
}

func (l *Logger) WarnContext(ctx context.Context, msg string) {
	l.logger.Warn(msg, extractFields(ctx))
}

func (l *Logger) ErrorContext(ctx context.Context, msg string) {
	l.logger.Error(msg, extractFields(ctx))
}

func Info(msg string) {
	getGlobalLogger().InfoContext(context.Background(), msg)
}

func Warn(msg string) {
	getGlobalLogger().WarnContext(context.Background(), msg)
}

func Error(msg string) {
	getGlobalLogger().ErrorContext(context.Background(), msg)
}

func InfoContext(ctx context.Context, msg string) {
	getGlobalLogger().InfoContext(ctx, msg)
}

func WarnContext(ctx context.Context, msg string) {
	getGlobalLogger().WarnContext(ctx, msg)
}

func ErrorContext(ctx context.Context, msg string) {
	getGlobalLogger().ErrorContext(ctx, msg)
}

type Field interface {
	Key() string
	Value() any
}

type simpleField struct {
	key string
	val any
}

func (f simpleField) Key() string { return f.key }
func (f simpleField) Value() any  { return f.val }

// NewField создает новое поле
func NewField(key string, value any) Field {
	return simpleField{key: key, val: value}
}

func fieldsToAny(fields []Field) []any {
	out := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		out = append(out, f.Key(), f.Value())
	}
	return out
}

type ctxKeyType struct{}

var ctxKey = ctxKeyType{}

func WithContext(ctx context.Context, fields ...Field) context.Context {
	existing, _ := ctx.Value(ctxKey).([]any)
	for _, f := range fields {
		existing = append(existing, f.Key(), f.Value())
	}
	return context.WithValue(ctx, ctxKey, existing)
}

func extractFields(ctx context.Context) []any {
	if ctx == nil {
		return nil
	}
	if fields, ok := ctx.Value(ctxKey).([]any); ok {
		return fields
	}
	return nil
}
