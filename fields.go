package alog

import (
	"strconv"
	"strings"
)

type Field interface {
	Key() string
	Render(ctx LogContext) string
}

type FieldOptions struct {
	Fields []Field
}

func Fields(fs ...Field) FieldOptions {
	return FieldOptions{Fields: fs}
}

/* ---------------- Time ---------------- */
type TimeField struct {
	Format string
}

func NewTimeField(format string) Field {
	return &TimeField{Format: format}
}

func (f *TimeField) Key() string { return "time" }

func (f *TimeField) Render(ctx LogContext) string {
	return ctx.Time.Format(f.Format)
}

/* ---------------- Level ---------------- */
type LevelField struct {
	Style string // upper / lower / title
}

func NewLevelField(style string) Field {
	return &LevelField{Style: style}
}

func (f *LevelField) Key() string { return "level" }

func (f *LevelField) Render(ctx LogContext) string {
	levelStr := ctx.Level.String()

	switch f.Style {
	case "lower":
		return strings.ToLower(levelStr)
	case "title":
		return strings.Title(strings.ToLower(levelStr))
	default:
		return strings.ToUpper(levelStr)
	}
}

/* ---------------- Msg ---------------- */
type MsgField struct{}

func NewMsgField() Field {
	return &MsgField{}
}

func (f *MsgField) Key() string { return "msg" }

func (f *MsgField) Render(ctx LogContext) string {
	return ctx.Msg
}

/* ---------------- File ---------------- */
type FileField struct{}

func NewFileField() Field {
	return &FileField{}
}

func (f *FileField) Key() string { return "file" }

func (f *FileField) Render(ctx LogContext) string {
	return ctx.File
}

/* ---------------- ShortFile ---------------- */
type ShortFileField struct{}

func NewShortFileField() Field {
	return &ShortFileField{}
}

func (f *ShortFileField) Key() string { return "short_file" }

func (f *ShortFileField) Render(ctx LogContext) string {
	return ctx.ShortFile
}

/* ---------------- Line ---------------- */
type LineField struct{}

func NewLineField() Field {
	return &LineField{}
}

func (f *LineField) Key() string { return "line" }

func (f *LineField) Render(ctx LogContext) string {
	return strconv.Itoa(ctx.Line)
}

/* ---------------- PID ---------------- */
type PIDField struct{}

func NewPIDField() Field {
	return &PIDField{}
}

func (f *PIDField) Key() string { return "pid" }

func (f *PIDField) Render(ctx LogContext) string {
	return strconv.Itoa(ctx.PID)
}

/* ---------------- GID ---------------- */
type GIDField struct{}

func NewGIDField() Field {
	return &GIDField{}
}

func (f *GIDField) Key() string { return "gid" }

func (f *GIDField) Render(ctx LogContext) string {
	return strconv.Itoa(ctx.GID)
}

/* ---------------- TimeStamp ---------------- */
type TimeStampField struct{}

func NewTimeStampField() Field {
	return &TimeStampField{}
}

func (f *TimeStampField) Key() string { return "time_stamp" }

func (f *TimeStampField) Render(ctx LogContext) string {
	return strconv.FormatInt(ctx.UnixNano, 10)
}

/* ---------------- TimeZone ---------------- */
type TimeZoneField struct{}

func NewTimeZoneField() Field {
	return &TimeZoneField{}
}

func (f *TimeZoneField) Key() string { return "time_zone" }

func (f *TimeZoneField) Render(ctx LogContext) string {
	return ctx.TZ
}

/* ---------------- TraceID ---------------- */
type TraceIDField struct{}

func NewTraceIDField() Field {
	return &TraceIDField{}
}

func (f *TraceIDField) Key() string { return "trace_id" }

func (f *TraceIDField) Render(ctx LogContext) string {
	return ctx.TraceID
}

/* ---------------- RequestID ---------------- */
type RequestIDField struct{}

func NewRequestIDField() Field {
	return &RequestIDField{}
}

func (f *RequestIDField) Key() string { return "request_id" }

func (f *RequestIDField) Render(ctx LogContext) string {
	return ctx.RequestID
}

/* ---------------- RawMsg ---------------- */
type RawMsgField struct{}

func NewRawMsgField() Field {
	return &RawMsgField{}
}

func (f *RawMsgField) Key() string { return "raw_msg" }

func (f *RawMsgField) Render(ctx LogContext) string {
	return ctx.RawMsg
}

/* ---------------- FileLine ---------------- */
type FileLineField struct {
	Prefix string
	Suffix string
}

func NewFileLineField(prefix, suffix string) Field {
	return &FileLineField{
		Prefix: prefix,
		Suffix: suffix,
	}
}

func (f *FileLineField) Key() string {
	return "fileline"
}

func (f *FileLineField) Render(ctx LogContext) string {
	v := ctx.ShortFile + ":" + strconv.Itoa(ctx.Line)
	return f.Prefix + v + f.Suffix
}

/* ---------------- Ext ---------------- */
type ExtField struct{}

func NewExtField() Field {
	return &ExtField{}
}

func (f *ExtField) Key() string { return "ext" }

func (f *ExtField) Render(ctx LogContext) string {
	if len(ctx.Ext) == 0 {
		return ""
	}

	var sb strings.Builder
	for k, v := range ctx.Ext {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
		sb.WriteString(" ")
	}
	return strings.TrimSpace(sb.String())
}
