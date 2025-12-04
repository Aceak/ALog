package alog

import (
	"fmt"
	"os"
)

type ConsoleSink struct{}

func NewConsoleSink() *ConsoleSink {
	return &ConsoleSink{}
}

func (s *ConsoleSink) Write(line string) {
	fmt.Fprintln(os.Stdout, line)
}
