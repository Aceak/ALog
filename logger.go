package alog

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Sink interface {
	Write(line string)
}

type Logger struct {
	level     Level
	formatter *Formatter
	sink      Sink
}

func NewLogger(level Level, formatter *Formatter, sink Sink) *Logger {
	return &Logger{
		level:     level,
		formatter: formatter,
		sink:      sink,
	}
}

func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}

	_, file, line, _ := runtime.Caller(2)

	now := time.Now()
	ctx := LogContext{
		Level:     level,
		Time:      now,
		UnixNano:  now.UnixNano(),
		TZ:        now.Location().String(),
		Msg:       msg,
		RawMsg:    msg,
		File:      file,
		ShortFile: filepath.Base(file),
		Line:      line,
		PID:       os.Getpid(),
		GID:       getGID(),
		Ext:       map[string]string{},
	}

	lineText := l.formatter.Format(ctx)
	l.sink.Write(lineText)

	if level == FATAL {
		panic("FATAL log encountered")
	}
}

func (l *Logger) Trace(msg string) { l.log(TRACE, msg) }
func (l *Logger) Debug(msg string) { l.log(DEBUG, msg) }
func (l *Logger) Info(msg string)  { l.log(INFO, msg) }
func (l *Logger) Warn(msg string)  { l.log(WARN, msg) }
func (l *Logger) Error(msg string) { l.log(ERROR, msg) }
func (l *Logger) Panic(msg string) { l.log(PANIC, msg) }
func (l *Logger) Fatal(msg string) { l.log(FATAL, msg) }
