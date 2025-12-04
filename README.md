# ALog 使用文档

## 1. 项目简介

ALog 是一个轻量级、高性能的 Go 语言日志库，提供了灵活的日志格式配置和多种输出方式，支持自定义日志字段和输出目标。

## 2. 安装方法

```bash
go get -u github.com/Aceak/ALog
```

## 3. 快速开始

### 3.1 全局日志 API

```go
package main

import (
    "github.com/Aceak/ALog"
)

func main() {
    alog.Trace("This is a trace message")
    alog.Debug("This is a debug message")
    alog.Info("This is an info message")
    alog.Warn("This is a warning message")
    alog.Error("This is an error message")
    // alog.Panic("This will cause panic")
    // alog.Fatal("This will cause fatal exit")
}
```

### 3.2 输出示例

```
2024-01-01 12:00:00 UTC INFO  [main.go:10] This is an info message
2024-01-01 12:00:00 UTC WARN  [main.go:11] This is a warning message
2024-01-01 12:00:00 UTC ERROR [main.go:12] This is an error message
```

## 4. 核心功能

### 4.1 自定义日志器

```go
package main

import (
    "github.com/Aceak/ALog"
)

func main() {
    // 创建自定义日志器
    logger := alog.NewLogger(
        alog.DEBUG,  // 设置日志级别为 DEBUG
        alog.NewFormatter(" | ",  // 设置字段分隔符
            alog.NewTimeField("2006-01-02 15:04:05"),  // 时间字段
            alog.NewLevelField("lower"),  // 日志级别字段，小写格式
            alog.NewFileLineField("", ""),  // 文件行号字段
            alog.NewMsgField(),  // 日志消息字段
        ),
        alog.NewConsoleSink(),  // 控制台输出
    )

    // 使用自定义日志器
    logger.Debug("This is a debug message")
    logger.Info("This is an info message")
}
```

### 4.2 自定义日志格式

ALog 支持通过字段组合来自定义日志格式，内置了多种日志字段：

| 字段类型 | 说明 | 示例 |
|---------|------|------|
| TimeField | 时间字段 | `NewTimeField("2006-01-02 15:04:05")` |
| LevelField | 日志级别字段 | `NewLevelField("upper")` |
| MsgField | 日志消息字段 | `NewMsgField()` |
| FileField | 文件路径字段 | `NewFileField()` |
| ShortFileField | 短文件名字段 | `NewShortFileField()` |
| LineField | 行号字段 | `NewLineField()` |
| FileLineField | 文件行号组合字段 | `NewFileLineField("[", "]")` |
| PIDField | 进程 ID 字段 | `NewPIDField()` |
| GIDField | Goroutine ID 字段 | `NewGIDField()` |
| TimeStampField | 时间戳字段 | `NewTimeStampField()` |
| TimeZoneField | 时区字段 | `NewTimeZoneField()` |
| TraceIDField | 追踪 ID 字段 | `NewTraceIDField()` |
| RequestIDField | 请求 ID 字段 | `NewRequestIDField()` |
| RawMsgField | 原始消息字段 | `NewRawMsgField()` |
| ExtField | 扩展字段 | `NewExtField()` |

### 4.3 日志级别

ALog 支持以下日志级别（按优先级从低到高）：

| 级别 | 说明 |
|------|------|
| TRACE | 追踪级别，用于详细的调试信息 |
| DEBUG | 调试级别，用于开发阶段的调试信息 |
| INFO | 信息级别，用于普通的运行时信息 |
| WARN | 警告级别，用于潜在的问题 |
| ERROR | 错误级别，用于错误信息 |
| PANIC | 恐慌级别，记录后会触发 panic |
| FATAL | 致命级别，记录后会触发程序退出 |

## 5. 高级功能

### 5.1 解析日志级别

```go
level := alog.ParseLevel("debug")  // 返回 alog.DEBUG
```

### 5.2 设置日志级别

```go
logger := alog.NewLogger(alog.INFO, formatter, sink)
// 日志级别可以通过 ParseLevel 动态设置
logger.SetLevel(alog.ParseLevel("debug"))
```

### 5.3 自定义输出目标

ALog 支持通过实现 Sink 接口来自定义输出目标：

```go
type MySink struct {
    // 自定义字段
}

func (s *MySink) Write(line string) {
    // 实现自定义输出逻辑
    fmt.Println("Custom sink:", line)
}

// 使用自定义 Sink
logger := alog.NewLogger(alog.INFO, formatter, &MySink{})
```

## 6. 最佳实践

### 6.1 在项目中使用

```go
package main

import (
    "github.com/Aceak/ALog"
)

// 初始化自定义日志器
var logger = alog.NewLogger(
    alog.Info,
    alog.NewFormatter(" | ",
        alog.NewTimeField("2006-01-02 15:04:05"),
        alog.NewLevelField("upper"),
        alog.NewFileLineField("[", "]"),
        alog.NewMsgField(),
    ),
    alog.NewConsoleSink(),
)

func main() {
    logger.Info("Application started")
    // 业务逻辑
    logger.Info("Application exited")
}
```

### 6.2 根据环境配置日志级别

```go
func init() {
    level := alog.Info
    if os.Getenv("ENV") == "development" {
        level = alog.Debug
    }
    logger = alog.NewLogger(level, formatter, sink)
}
```

## 7. 性能考虑

- ALog 使用了高效的字符串构建器 `strings.Builder` 来格式化日志
- 日志级别检查在最开始进行，低于当前级别的日志会被快速跳过
- 支持异步输出（通过自定义 Sink 实现）

## 8. 总结

ALog 是一个简单易用、功能丰富的 Go 语言日志库，提供了灵活的配置选项和良好的扩展性。它适合各种规模的 Go 项目使用，从简单的命令行工具到复杂的分布式系统。

通过合理配置日志格式和级别，可以帮助开发者更好地理解系统运行状态，快速定位和解决问题。