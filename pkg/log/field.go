package log

// Field представляет поле логирования (ключ-значение)
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
