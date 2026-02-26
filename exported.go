package alog

import "sync/atomic"

var std atomic.Value // stores *Logger

// setStd 设置全局 Logger，线程安全
func setStd(l *Logger) {
	std.Store(l)
}

// getStd 获取全局 Logger，线程安全
func getStd() *Logger {
	if v := std.Load(); v != nil {
		return v.(*Logger)
	}
	// 理论上 init() 会先 setStd，这里做兜底避免 nil
	l := NewLogger(
		INFO,
		NewFormatter(" ",
			NewTimeField("2006-01-02 15:04:05 MST"),
			NewLevelField("upper"),
			NewFileLineField("[", "]"),
			NewMsgField(),
		),
		NewConsoleSink(),
	)
	setStd(l)
	return l
}

// init 初始化全局 Logger，默认配置
func init() {
	_ = Init(DefaultConfig())
}

/* 全局 API */
func Trace(msg string) { getStd().Trace(msg) }
func Debug(msg string) { getStd().Debug(msg) }
func Info(msg string)  { getStd().Info(msg) }
func Warn(msg string)  { getStd().Warn(msg) }
func Error(msg string) { getStd().Error(msg) }
func Panic(msg string) { getStd().Panic(msg) }
func Fatal(msg string) { getStd().Fatal(msg) }

// SetFormatter 设置全局 Logger 的 Formatter
func SetFormatter(f *Formatter) {
	if f == nil {
		return
	}
	old := getStd()
	nl := cloneLogger(old)
	nl.formatter = f
	setStd(nl)
}
