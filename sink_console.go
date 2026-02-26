package alog

import (
	"os"
)

// ConsoleSink 实现基于控制台的日志输出
type ConsoleSink struct {
}

// NewConsoleSink 创建新的控制台日志 sink
func NewConsoleSink() *ConsoleSink {
	return &ConsoleSink{}
}

// Write 写入日志行
// line: 日志内容
func (s *ConsoleSink) Write(line string) {
	os.Stdout.WriteString(line + "\n")
}
