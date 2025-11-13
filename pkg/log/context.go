package log

import "context"

type ctxKeyType struct{}

var ctxKey = ctxKeyType{}

func WithContext(ctx context.Context, fields ...Field) context.Context {
	existing, _ := ctx.Value(ctxKey).([]Field)
	existing = append(existing, fields...)
	return context.WithValue(ctx, ctxKey, existing)
}

func extractFields(ctx context.Context) []any {
	var out []any
	if ctx == nil {
		return out
	}
	if fields, ok := ctx.Value(ctxKey).([]Field); ok {
		for _, f := range fields {
			out = append(out, f.Key(), f.Value())
		}
	}
	return out
}
