package alog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// RollingPolicy 定义日志轮转策略
type RollingPolicy int

const (
	RollNone   RollingPolicy = iota // 不轮转
	RollByDay                       // 按天轮转
	RollBySize                      // 按大小轮转
)

// FileOption 定义文件日志配置选项函数类型
type FileOption func(*FileSink)

// WithDayRolling 启用按天轮转策略
func WithDayRolling() FileOption {
	return func(s *FileSink) {
		s.policy = RollByDay
		s.currentDay = time.Now().Format("2006-01-02")
	}
}

// WithSizeRolling 启用按大小轮转策略
// size 可以是数字（字节）或字符串（如 "10MB"）
func WithSizeRolling(size interface{}) FileOption {
	return func(s *FileSink) {
		parsedSize := parseSize(size)
		if parsedSize <= 0 {
			// 使用默认值 10MB，不输出警告
			parsedSize = 10 * 1024 * 1024 // 10MB default
		}
		s.policySize = RollBySize
		s.maxSize = parsedSize
	}
}

// WithMaxDays 设置日志文件的最大保留天数
// 0 表示无限制
func WithMaxDays(days int) FileOption {
	return func(s *FileSink) {
		if days < 0 {
			// 使用默认值 0（无限制），不输出警告
			days = 0
		}
		s.maxDays = days
	}
}

// WithMaxArchives 设置日志文件的最大归档数量
// 0 表示无限制
func WithMaxArchives(archives int) FileOption {
	return func(s *FileSink) {
		if archives < 0 {
			// 使用默认值 0（无限制），不输出警告
			archives = 0
		}
		s.maxArchives = archives
	}
}

// WithFileMode 设置日志文件的权限模式
func WithFileMode(mode os.FileMode) FileOption {
	return func(s *FileSink) {
		s.fileMode = mode
	}
}

// WithDirMode 设置日志目录的权限模式
func WithDirMode(mode os.FileMode) FileOption {
	return func(s *FileSink) {
		s.dirMode = mode
	}
}

// FileSink 实现基于文件的日志输出
// 支持按天和按大小轮转，以及自动清理旧日志

type FileSink struct {
	mu sync.Mutex // 互斥锁，用于并发安全

	file *os.File // 当前打开的日志文件
	path string   // 日志文件路径

	policy     RollingPolicy // 按天轮转策略
	policySize RollingPolicy // 按大小轮转策略
	maxSize    int64         // 按大小轮转的最大文件大小

	currentDay string // 当前日期，格式：2006-01-02

	maxDays     int         // 日志文件的最大保留天数
	maxArchives int         // 日志文件的最大归档数量
	fileMode    os.FileMode // 文件权限模式
	dirMode     os.FileMode // 目录权限模式
}

// NewFileSink 创建新的文件日志 sink
// path: 日志文件路径
// opts: 配置选项
func NewFileSink(path string, opts ...FileOption) (*FileSink, error) {
	sink := &FileSink{
		path:     path,
		fileMode: 0644, // 默认文件权限
		dirMode:  0755, // 默认目录权限
	}

	for _, opt := range opts {
		opt(sink)
	}

	if err := sink.openNewFile(); err != nil {
		return nil, err
	}

	return sink, nil
}

// openNewFile 打开新的日志文件
// 自动创建不存在的目录
func (s *FileSink) openNewFile() error {
	// 确保目录存在
	dir := filepath.Dir(s.path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, s.dirMode); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, s.fileMode)
	if err != nil {
		return err
	}
	s.file = file
	return nil
}

// Write 写入日志行
// line: 日志内容
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

// checkDayRolling 检查是否需要按天轮转
// t: 当前时间
func (s *FileSink) checkDayRolling(t time.Time) {
	day := t.Format("2006-01-02")
	if day == s.currentDay {
		return
	}

	yesterday := s.currentDay
	s.currentDay = day

	// 忽略错误，因为日志库不应该产生额外的输出
	_ = s.file.Close()

	s.rotateFinal(yesterday)

	s.cleanup()
	// 忽略错误，因为日志库不应该产生额外的输出
	_ = s.openNewFile()
}

// checkSizeRolling 检查是否需要按大小轮转
// line: 要写入的日志内容
func (s *FileSink) checkSizeRolling(line string) {
	if s.maxSize <= 0 {
		return
	}

	info, err := s.file.Stat()
	if err != nil {
		// 忽略错误，因为日志库不应该产生额外的输出
		return
	}

	if info.Size() >= s.maxSize {
		s.rotateSize()
	}
}

// rotateSize 执行按大小轮转
func (s *FileSink) rotateSize() {
	// 忽略错误，因为日志库不应该产生额外的输出
	_ = s.file.Close()

	i := 1
	for {
		old := fmt.Sprintf("%s.%d", s.path, i)
		if _, err := os.Stat(old); os.IsNotExist(err) {
			break
		}
		i++
	}

	for j := i - 1; j > 0; j-- {
		old := fmt.Sprintf("%s.%d", s.path, j)
		new := fmt.Sprintf("%s.%d", s.path, j+1)
		// 忽略错误，因为日志库不应该产生额外的输出
		_ = os.Rename(old, new)
	}

	// 忽略错误，因为日志库不应该产生额外的输出
	_ = os.Rename(s.path, fmt.Sprintf("%s.1", s.path))
	// 忽略错误，因为日志库不应该产生额外的输出
	_ = s.openNewFile()
}

// rotateFinal 执行按天轮转的最终操作
// day: 要归档的日期
func (s *FileSink) rotateFinal(day string) {
	archive := fmt.Sprintf("%s.%s", s.path, day)

	if _, err := os.Stat(archive); err == nil {
		i := 1
		for {
			candidate := fmt.Sprintf("%s.%s.%d", s.path, day, i)
			if _, err := os.Stat(candidate); os.IsNotExist(err) {
				archive = candidate
				break
			}
			i++
		}
	}

	// 忽略错误，因为日志库不应该产生额外的输出
	_ = os.Rename(s.path, archive)
}

// cleanup 清理旧日志文件
func (s *FileSink) cleanup() {
	if s.maxDays <= 0 && s.maxArchives <= 0 {
		return
	}

	dir := "."
	base := s.path

	if idx := len(s.path) - len(s.path[strings.LastIndex(s.path, "/")+1:]); idx > 0 {
		dir = s.path[:strings.LastIndex(s.path, "/")]
		base = s.path[strings.LastIndex(s.path, "/")+1:]
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		// 忽略错误，因为日志库不应该产生额外的输出
		return
	}

	type fileInfo struct {
		name string
		mod  time.Time
	}

	var archives []fileInfo

	prefix := base + "."

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			// 忽略错误，因为日志库不应该产生额外的输出
			continue
		}
		archives = append(archives, fileInfo{
			name: name,
			mod:  info.ModTime(),
		})
	}

	sort.Slice(archives, func(i, j int) bool {
		return archives[i].mod.After(archives[j].mod)
	})

	if s.maxArchives > 0 && len(archives) > s.maxArchives {
		for _, f := range archives[s.maxArchives:] {
			// 忽略错误，因为日志库不应该产生额外的输出
			_ = os.Remove(filepath.Join(dir, f.name))
		}
		archives = archives[:s.maxArchives]
	}

	if s.maxDays > 0 {
		expireBefore := time.Now().AddDate(0, 0, -s.maxDays)
		for _, f := range archives {
			if f.mod.Before(expireBefore) {
				// 忽略错误，因为日志库不应该产生额外的输出
				_ = os.Remove(filepath.Join(dir, f.name))
			}
		}
	}
}
