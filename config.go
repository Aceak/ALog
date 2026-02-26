package alog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// LoggerConfig 配置：用于装配全局 std logger 或构建新 logger
type LoggerConfig struct {
	Level  string
	Format string

	EnableConsole bool
	FilePath      string

	PanicBehavior string
	FatalBehavior string
}

func DefaultConfig() LoggerConfig {
	return LoggerConfig{
		Level:         "debug",
		Format:        "[{time}] [{level}] {msg}",
		EnableConsole: true,
		FilePath:      "",
		PanicBehavior: "panic",
		FatalBehavior: "exit",
	}
}

// Init：用配置替换全局 std logger
func Init(cfg LoggerConfig) error {
	l, err := NewLoggerFromConfig(cfg)
	if err != nil {
		return err
	}
	setStd(l)
	return nil
}

// NewLoggerFromConfig：根据配置创建 Logger
func NewLoggerFromConfig(cfg LoggerConfig) (*Logger, error) {
	level := ParseLevel(strings.ToLower(strings.TrimSpace(cfg.Level)))

	formatter, err := NewFormatterFromTemplate(cfg.Format)
	if err != nil {
		return nil, err
	}

	var sinks []Sink
	if cfg.EnableConsole {
		sinks = append(sinks, NewConsoleSink())
	}
	if filePath := strings.TrimSpace(cfg.FilePath); filePath != "" {
		// 验证文件路径
		dir := filepath.Dir(filePath)
		if dir != "" && dir != "." {
			// 检查目录是否存在，不存在则尝试创建
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create log directory: %v", err)
			}
		}

		fs, err := NewFileSink(filePath, WithDayRolling())
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, fs)
	}

	var sink Sink
	switch len(sinks) {
	case 0:
		sink = NewConsoleSink()
	case 1:
		sink = sinks[0]
	default:
		sink = NewMultiSink(sinks...)
	}

	l := NewLogger(level, formatter, sink)

	// 行为策略接入
	switch strings.ToLower(strings.TrimSpace(cfg.PanicBehavior)) {
	case "none":
		l.SetPanicBehavior(PanicNone)
	default:
		l.SetPanicBehavior(PanicPanic)
	}

	switch strings.ToLower(strings.TrimSpace(cfg.FatalBehavior)) {
	case "none":
		l.SetFatalBehavior(FatalNone)
	case "panic":
		l.SetFatalBehavior(FatalPanic)
	default:
		l.SetFatalBehavior(FatalExit)
	}

	return l, nil
}

type TextField struct{ Text string }

func (f *TextField) Key() string                  { return "text" }
func (f *TextField) Render(ctx LogContext) string { return f.Text }

var (
	templateCache = make(map[string][]Field)
	cacheMutex    sync.RWMutex
)

func NewFormatterFromTemplate(tpl string) (*Formatter, error) {
	tpl = strings.TrimSpace(tpl)
	if tpl == "" {
		tpl = "[{time}] [{level}] {msg}"
	}

	// 尝试从缓存获取
	cacheMutex.RLock()
	cached, ok := templateCache[tpl]
	cacheMutex.RUnlock()

	if ok {
		// 防御式拷贝：避免外部修改影响缓存
		fields := append([]Field(nil), cached...)
		return NewFormatter("", fields...), nil
	}

	fields, err := parseTemplateToFields(tpl)
	if err != nil {
		return nil, err
	}

	toCache := append([]Field(nil), fields...)
	cacheMutex.Lock()
	templateCache[tpl] = toCache
	cacheMutex.Unlock()

	out := append([]Field(nil), toCache...)
	return NewFormatter("", out...), nil
}

func parseTemplateToFields(tpl string) ([]Field, error) {
	var out []Field
	var textBuf strings.Builder

	flushText := func() {
		if textBuf.Len() > 0 {
			out = append(out, &TextField{Text: textBuf.String()})
			textBuf.Reset()
		}
	}

	for i := 0; i < len(tpl); i++ {
		ch := tpl[i]

		// 转义：\{ \} \\
		if ch == '\\' && i+1 < len(tpl) {
			next := tpl[i+1]
			if next == '{' || next == '}' || next == '\\' {
				textBuf.WriteByte(next)
				i++
				continue
			}
			textBuf.WriteByte(ch)
			continue
		}

		if ch != '{' {
			textBuf.WriteByte(ch)
			continue
		}

		// 遇到 {：先把之前的纯文本刷出去
		flushText()

		// 线性找 }
		j := i + 1
		for ; j < len(tpl); j++ {
			if tpl[j] == '}' {
				break
			}
			// 注意：模板内部不支持转义占位符内容，这里只找第一个 }
		}

		if j >= len(tpl) || tpl[j] != '}' {
			return nil, fmt.Errorf("format: missing '}'")
		}

		rawToken := tpl[i+1 : j]
		token := strings.TrimSpace(rawToken)
		if token == "" {
			return nil, fmt.Errorf("format: empty placeholder at pos %d, template=%q", i, tpl)
		}

		f, err := fieldFromToken(token)
		if err != nil {
			return nil, fmt.Errorf("format: bad placeholder %q at pos %d, template=%q: %w", token, i, tpl, err)
		}
		out = append(out, f)

		i = j // 跳过 }
	}

	flushText()
	return out, nil
}

func fieldFromToken(token string) (Field, error) {
	name := token
	arg := ""

	if idx := strings.IndexByte(token, ':'); idx >= 0 {
		name = strings.TrimSpace(token[:idx])
		arg = strings.TrimSpace(token[idx+1:])
	}

	switch name {
	case "time":
		layout := arg
		if layout == "" {
			layout = "2006-01-02 15:04:05 MST"
		}
		return NewTimeField(layout), nil

	case "level":
		style := arg
		if style == "" {
			style = "upper"
		}
		return NewLevelField(style), nil

	case "msg":
		return NewMsgField(), nil

	case "file":
		return NewFileField(), nil
	case "short_file":
		return NewShortFileField(), nil
	case "line":
		return NewLineField(), nil

	case "pid":
		return NewPIDField(), nil
	case "gid":
		return NewGIDField(), nil

	case "time_stamp":
		return NewTimeStampField(), nil
	case "time_zone":
		return NewTimeZoneField(), nil

	case "trace_id":
		return NewTraceIDField(), nil
	case "request_id":
		return NewRequestIDField(), nil

	case "raw_msg":
		return NewRawMsgField(), nil
	case "ext":
		return NewExtField(), nil

	case "fileline":
		if arg == "" {
			return NewFileLineField("[", "]"), nil
		}
		parts := strings.SplitN(arg, ",", 2)
		prefix := parts[0]
		suffix := ""
		if len(parts) == 2 {
			suffix = parts[1]
		}
		return NewFileLineField(prefix, suffix), nil

	default:
		return nil, fmt.Errorf("unknown placeholder: %s", token)
	}
}
