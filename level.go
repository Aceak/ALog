package alog

type Level int

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

func ParseLevel(s string) Level {
	switch s {
	case "trace":
		return TRACE
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	case "panic":
		return PANIC
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

func (l Level) String() string {
	switch l {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO "
	case WARN:
		return "WARN "
	case ERROR:
		return "ERROR"
	case PANIC:
		return "PANIC"
	case FATAL:
		return "FATAL"
	default:
		return "INFO "
	}
}
