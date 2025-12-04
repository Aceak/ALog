package alog

import (
	"strings"
)

type Formatter struct {
	fields []Field
	sep    string
}

func NewFormatter(sep string, fields ...Field) *Formatter {
	return &Formatter{
		fields: fields,
		sep:    sep,
	}
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
