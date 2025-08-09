# GoSpring 日志系统

GoSpring 框架提供了一个强大而灵活的日志系统，用于记录框架内部的各种事件，包括容器创建、组件注册、依赖注入、生命周期管理等。

## 设计理念

GoSpring 的日志系统采用了事件驱动的设计，参考了 Uber 的 fx 框架的日志设计模式：

- **事件驱动**：所有日志都以事件的形式记录
- **可扩展**：支持自定义日志器实现
- **类型安全**：使用强类型的事件结构
- **灵活配置**：支持多种日志器组合

## 核心组件

### 1. Event 接口

所有日志事件都实现 `Event` 接口：

```go
type Event interface {
    String() string
}
```

### 2. Logger 接口

日志器接口定义了日志记录的核心方法：

```go
type Logger interface {
    LogEvent(event Event)
}
```

### 3. 内置日志器

#### NopLogger
空操作日志器，不记录任何日志：

```go
logger := logging.NewNopLogger()
```

#### ConsoleLogger
控制台日志器，将日志输出到标准输出：

```go
logger := logging.NewConsoleLogger()
```

#### StandardLogger
标准库日志器，使用 Go 标准库的 log 包：

```go
import "log"
import "os"

stdLog := log.New(os.Stdout, "[GoSpring] ", log.LstdFlags)
logger := logging.NewStandardLogger(stdLog)
```

#### MultiLogger
多重日志器，同时向多个日志器输出：

```go
logger1 := logging.NewConsoleLogger()
logger2 := logging.NewStandardLogger(stdLog)
multiLogger := logging.NewMultiLogger(logger1, logger2)
```

#### FilteredLogger
过滤日志器，只记录满足条件的事件：

```go
baseLogger := logging.NewConsoleLogger()
filteredLogger := logging.NewFilteredLogger(baseLogger, func(event logging.Event) bool {
    // 只记录错误事件
    switch event.(type) {
    case *logging.DependencyInjectionFailed:
        return true
    default:
        return false
    }
})
```

#### LeveledLogger
分级日志器，支持不同级别的日志记录：

```go
baseLogger := logging.NewConsoleLogger()
leveledLogger := logging.NewLeveledLogger(baseLogger, logging.LogLevelInfo)
```

## 事件类型

### 容器事件

- **ContainerCreated**: 容器创建事件
- **ComponentRegistered**: 组件注册事件
- **ComponentCreated**: 组件创建事件
- **ComponentDestroyed**: 组件销毁事件

### 依赖注入事件

- **DependencyInjected**: 依赖注入成功事件
- **DependencyInjectionFailed**: 依赖注入失败事件

### 生命周期事件

- **LifecycleStarting**: 生命周期方法开始执行事件
- **LifecycleStarted**: 生命周期方法执行完成事件
- **LifecycleStopping**: 生命周期停止开始事件
- **LifecycleStopped**: 生命周期停止完成事件

### 应用上下文事件

- **ContextStarting**: 应用上下文启动开始事件
- **ContextStarted**: 应用上下文启动完成事件
- **ContextStopping**: 应用上下文停止开始事件
- **ContextStopped**: 应用上下文停止完成事件

### 扫描事件

- **ScanStarting**: 组件扫描开始事件
- **ScanCompleted**: 组件扫描完成事件

## 使用方法

### 1. 基本使用

```go
import (
    "gospring/context"
    "gospring/logging"
)

// 创建带有日志器的应用上下文
logger := logging.NewConsoleLogger()
ctx := context.NewApplicationContextWithLogger(logger)

// 注册组件
ctx.RegisterComponent(&UserService{})

// 启动应用上下文（会记录相关日志）
ctx.Start()
```

### 2. 自定义日志器

```go
// 创建自定义过滤器
customLogger := logging.NewFilteredLogger(
    logging.NewConsoleLogger(),
    func(event logging.Event) bool {
        // 只记录组件相关事件
        switch event.(type) {
        case *logging.ComponentCreated, *logging.ComponentDestroyed:
            return true
        default:
            return false
        }
    },
)

ctx := context.NewApplicationContextWithLogger(customLogger)
```

### 3. 多重日志输出

```go
// 同时输出到控制台和文件
consoleLogger := logging.NewConsoleLogger()
fileLogger := logging.NewStandardLogger(log.New(file, "", log.LstdFlags))
multiLogger := logging.NewMultiLogger(consoleLogger, fileLogger)

ctx := context.NewApplicationContextWithLogger(multiLogger)
```

### 4. 运行时更改日志器

```go
ctx := context.NewApplicationContext()

// 运行时设置日志器
ctx.SetLogger(logging.NewConsoleLogger())

// 获取当前日志器
currentLogger := ctx.GetLogger()
```

## 最佳实践

1. **生产环境**：使用 `FilteredLogger` 只记录重要事件，避免日志过多
2. **开发环境**：使用 `ConsoleLogger` 查看详细的框架运行信息
3. **调试模式**：使用 `MultiLogger` 同时输出到控制台和文件
4. **性能敏感**：使用 `NopLogger` 完全禁用日志记录

## 示例

完整的使用示例请参考 `examples/logging/main.go` 文件。

## 扩展

你可以通过实现 `Logger` 接口来创建自定义的日志器：

```go
type CustomLogger struct {
    // 自定义字段
}

func (l *CustomLogger) LogEvent(event logging.Event) {
    // 自定义日志处理逻辑
    switch e := event.(type) {
    case *logging.ComponentCreated:
        // 处理组件创建事件
    case *logging.DependencyInjected:
        // 处理依赖注入事件
    default:
        // 处理其他事件
    }
}
```