package alog

import (
	"strings"
)

type Formatter struct {
	fields []Field
	sep    string
}

func NewFormatter(sep string, fields ...Field) *Formatter {
	return NewFormatterWithOptions(sep, Fields(fields...))
}

func NewFormatterWithOptions(sep string, opt FieldOptions) *Formatter {
	fs := make([]Field, 0, len(opt.Fields))
	for _, it := range opt.Fields {
		if it == nil {
			continue
		}
		fs = append(fs, it)
	}
	return &Formatter{fields: fs, sep: sep}
}

func (f *Formatter) SetSeparator(sep string) {
	f.sep = sep
}

func (f *Formatter) Format(ctx LogContext) string {
	var sb strings.Builder

	first := true

	for _, field := range f.fields {
		val := field.Render(ctx)
		if val == "" {
			continue
		}

		if !first {
			sb.WriteString(f.sep)
		}

		sb.WriteString(val)
		first = false
	}

	return sb.String()
}
