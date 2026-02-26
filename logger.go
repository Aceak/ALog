package alog

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 缓存当前包的路径，用于判断是否为 alog 包内部调用
var alogPackagePath string

// 包初始化函数，获取当前包的路径
func init() {
	// 获取当前包的路径
	_, file, _, _ := runtime.Caller(0)
	alogPackagePath = file
	// 找到最后一个路径分隔符，截取到包目录
	if idx := strings.LastIndexAny(alogPackagePath, "/\\"); idx >= 0 {
		alogPackagePath = alogPackagePath[:idx]
	}
}

type PanicBehavior string
type FatalBehavior string

const (
	PanicPanic PanicBehavior = "panic"
	PanicNone  PanicBehavior = "none"

	FatalExit  FatalBehavior = "exit"
	FatalPanic FatalBehavior = "panic"
	FatalNone  FatalBehavior = "none"
)

type Logger struct {
	level     Level
	formatter *Formatter
	sink      Sink

	panicBehavior PanicBehavior
	fatalBehavior FatalBehavior

	traceID   string
	requestID string
	ext       map[string]string
}

func NewLogger(level Level, formatter *Formatter, sink Sink) *Logger {
	return &Logger{
		level:     level,
		formatter: formatter,
		sink:      sink,

		panicBehavior: PanicPanic,
		fatalBehavior: FatalExit,
	}
}

func (l *Logger) SetPanicBehavior(behavior PanicBehavior) { l.panicBehavior = behavior }
func (l *Logger) SetFatalBehavior(behavior FatalBehavior) { l.fatalBehavior = behavior }

func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}

	file, line := findCaller(4)

	now := time.Now()

	// 准备扩展字段
	ext := map[string]string{}
	if len(l.ext) > 0 {
		ext = make(map[string]string, len(l.ext))
		for k, v := range l.ext {
			ext[k] = v
		}
	}

	// 添加 traceID 和 requestID（如果有）
	if l.traceID != "" {
		ext["trace_id"] = l.traceID
	}
	if l.requestID != "" {
		ext["request_id"] = l.requestID
	}

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
		Ext:       ext,
	}

	lineText := l.formatter.Format(ctx)
	l.sink.Write(lineText)

	if level == PANIC {
		switch l.panicBehavior {
		case PanicPanic:
			panic(msg)
		case PanicNone:
			return
		}
	}

	if level == FATAL {
		switch l.fatalBehavior {
		case FatalExit:
			os.Exit(1)
		case FatalPanic:
			panic(msg)
		case FatalNone:
			return
		}
	}
}

// findCaller 返回第一帧“不是 alog 包内部”的调用位置。
// skipBase 用来跳过 runtime.Callers / findCaller 自身等固定栈帧。
func findCaller(skipBase int) (string, int) {
	// 调整栈大小，根据实际需要设置
	const maxFrames = 64 // 通常 16 个栈帧足够
	pcs := make([]uintptr, maxFrames)
	n := runtime.Callers(skipBase, pcs)
	if n == 0 {
		return "unknown", 0
	}

	frames := runtime.CallersFrames(pcs[:n])

	prefix := alogPackagePath + string(os.PathSeparator)
	for {
		frame, more := frames.Next()

		// 使用更准确的路径判断：检查文件路径是否在 alog 包路径下
		if !strings.HasPrefix(frame.File, prefix) {
			return frame.File, frame.Line
		}

		if !more {
			break
		}
	}

	return "unknown", 0
}

// cloneLogger 深拷贝 Logger，避免共享 map
func cloneLogger(l *Logger) *Logger {
	if l == nil {
		return nil
	}
	nl := *l // 浅拷贝结构体

	// ext 深拷贝，避免共享 map
	if l.ext != nil {
		nl.ext = make(map[string]string, len(l.ext))
		for k, v := range l.ext {
			nl.ext[k] = v
		}
	}
	return &nl
}

// WithTraceID 返回一个携带 traceID 的 logger（
func (l *Logger) WithTraceID(id string) *Logger {
	nl := cloneLogger(l)
	nl.traceID = id
	return nl
}

// WithRequestID 返回一个携带 requestID 的 logger
func (l *Logger) WithRequestID(id string) *Logger {
	nl := cloneLogger(l)
	nl.requestID = id
	return nl
}

// WithExt 设置一个扩展字段（k=v）
func (l *Logger) WithExt(k, v string) *Logger {
	if k == "" {
		return l
	}
	nl := cloneLogger(l)
	if nl.ext == nil {
		nl.ext = map[string]string{}
	}
	nl.ext[k] = v
	return nl
}

// WithExtMap 批量设置扩展字段。nil/空 map 直接返回克隆或原对象都行，这里返回克隆以保持一致性。
func (l *Logger) WithExtMap(m map[string]string) *Logger {
	nl := cloneLogger(l)
	if len(m) == 0 {
		return nl
	}
	if nl.ext == nil {
		nl.ext = map[string]string{}
	}
	for k, v := range m {
		if k == "" {
			continue
		}
		nl.ext[k] = v
	}
	return nl
}

func (l *Logger) Trace(msg string) { l.log(TRACE, msg) }
func (l *Logger) Debug(msg string) { l.log(DEBUG, msg) }
func (l *Logger) Info(msg string)  { l.log(INFO, msg) }
func (l *Logger) Warn(msg string)  { l.log(WARN, msg) }
func (l *Logger) Error(msg string) { l.log(ERROR, msg) }
func (l *Logger) Panic(msg string) { l.log(PANIC, msg) }
func (l *Logger) Fatal(msg string) { l.log(FATAL, msg) }

func WithTraceID(id string) *Logger   { return getStd().WithTraceID(id) }
func WithRequestID(id string) *Logger { return getStd().WithRequestID(id) }
func WithExt(k, v string) *Logger     { return getStd().WithExt(k, v) }
func WithExtMap(m map[string]string) *Logger {
	return getStd().WithExtMap(m)
}
