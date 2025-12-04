package alog

import "time"

type LogContext struct {
	Level Level
	Msg   string
	Time  time.Time

	File      string
	ShortFile string
	Line      int
	PID       int
	GID       int

	UnixNano int64
	TZ       string

	TraceID   string
	RequestID string
	RawMsg    string
	Ext       map[string]string
}
