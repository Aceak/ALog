package alog

import (
	"os"
	"sync"
	"time"
)

type RollingPolicy int

const (
	RollNone RollingPolicy = iota
	RollByDay
	RollBySize
)

type FileOption func(*FileSink)

func WithDayRolling() FileOption {
	return func(s *FileSink) {
		s.policy = RollByDay
		s.currentDay = time.Now().Format("2006-01-02")
	}
}

func WithSizeRolling(size interface{}) FileOption {
	return func(s *FileSink) {
		s.policySize = RollBySize
		s.maxSize = parseSize(size)
	}
}

func WithMaxDays(days int) FileOption {
	return func(s *FileSink) {
		s.maxDays = days
	}
}

func WithMaxArchives(archives int) FileOption {
	return func(s *FileSink) {
		s.maxArchives = archives
	}
}

type FileSink struct {
	mu sync.Mutex

	file *os.File
	path string

	policy     RollingPolicy
	policySize RollingPolicy
	maxSize    int64

	currentDay string

	maxDays     int
	maxArchives int
}

func NewFileSink(path string, opts ...FileOption) (*FileSink, error) {
	sink := &FileSink{
		path: path,
	}

	for _, opt := range opts {
		opt(sink)
	}

	if err := sink.openNewFile(); err != nil {
		return nil, err
	}

	return sink, nil
}

func (s *FileSink) openNewFile() error {
	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	s.file = file
	return nil
}

func (s *FileSink) Write(line string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	if s.policy == RollByDay {
		s.checkDayRolling(now)
	}

	if s.policySize == RollBySize {
		s.checkSizeRolling(line)
	}

	s.file.WriteString(line + "\n")
}

func (s *FileSink) checkDayRolling(t time.Time) {
	day := t.Format("2006-01-02")
	if day == s.currentDay {
		return
	}
}
