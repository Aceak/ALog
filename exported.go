package alog

var std *Logger

func init() {
	std = NewLogger(
		INFO,
		NewFormatter(" ",
			NewTimeField("2006-01-02 15:04:05 MST"),
			NewLevelField("upper"),
			NewFileLineField("[", "]"),
			NewMsgField(),
		),
		NewConsoleSink(),
	)
}

/* 全局 API */

func Trace(msg string) { std.Trace(msg) }
func Debug(msg string) { std.Debug(msg) }
func Info(msg string)  { std.Info(msg) }
func Warn(msg string)  { std.Warn(msg) }
func Error(msg string) { std.Error(msg) }
func Panic(msg string) { std.Panic(msg) }
func Fatal(msg string) { std.Fatal(msg) }
