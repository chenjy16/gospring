# GoSpring - Go IoC/DI Container Framework

[English](#english) | [中文](#chinese)

## English

A lightweight dependency injection framework for Go, inspired by Spring's IoC/DI container implementation, using Go's reflection, interfaces, and struct tags.

### ✨ Features

- 🚀 **Automatic Dependency Injection** - Zero-configuration dependency injection based on reflection
- 🏷️ **Tag-Driven** - Configuration using Go struct tags
- 🔄 **Lifecycle Management** - Complete Bean initialization and destruction lifecycle
- 📦 **Component Scanning** - Automatic component discovery and registration
- 🎯 **Interface Binding** - Automatic binding between interfaces and implementations
- 🛠️ **Multiple Scopes** - Support for singleton and prototype patterns
- 🧵 **Thread-Safe** - Concurrent-safe container implementation
- 🔧 **Type-Safe** - Safety guarantees based on Go's strong type system
- 📝 **Event Logging** - Comprehensive framework event recording and monitoring

### 📁 Project Structure

```
gospring/
├── container/          # Core container implementation
├── context/           # Application context
├── scanner/           # Component scanner
├── lifecycle/         # Lifecycle management
├── logging/           # Logging system
├── annotations/       # Annotation and tag processing
├── examples/          # Example code
│   ├── basic/        # Basic usage examples
│   ├── web/          # Web application examples
│   └── logging/      # Logging system examples
├── tests/            # Unit tests
├── docs/             # Documentation
└── README.md
```

### Quick Start

```go
package main

import (
    "fmt"
    "gospring/container"
)

// Define service interface
type UserService interface {
    GetUser(id int) string
}

// Implement service
type UserServiceImpl struct {
    Repository UserRepository `inject:""`
}

func (u *UserServiceImpl) GetUser(id int) string {
    return u.Repository.FindById(id)
}

// Define repository interface
type UserRepository interface {
    FindById(id int) string
}

// Implement repository
type UserRepositoryImpl struct{}

func (u *UserRepositoryImpl) FindById(id int) string {
    return fmt.Sprintf("User-%d", id)
}

func main() {
    // Create container
    c := container.NewContainer()
    
    // Register components
    c.RegisterSingleton("userRepository", &UserRepositoryImpl{})
    c.RegisterSingleton("userService", &UserServiceImpl{})
    
    // Get service
    userService := c.GetBean("userService").(UserService)
    result := userService.GetUser(1)
    fmt.Println(result) // Output: User-1
}
```

### Core Concepts

#### 1. Container
Manages the lifecycle and dependencies of all components.

#### 2. Dependency Injection
Automatically injects dependencies through `inject` tags.

#### 3. Component Scanning
Automatically discovers and registers components with specific tags.

#### 4. Lifecycle Management
Supports component initialization and destruction callbacks.

#### 5. Logging System
Provides event-driven logging with multiple logger implementations and flexible configuration.

### 🏷️ Tag Reference

| Tag | Description | Example |
|-----|-------------|---------|
| `inject:""` | Mark field for injection | `Repository UserRepo \`inject:""\`` |
| `inject:"beanName"` | Specify Bean name for injection | `Cache CacheService \`inject:"redisCache"\`` |
| `component:""` | Mark as component | `_ string \`component:"userService"\`` |
| `singleton:"true"` | Mark as singleton | `_ string \`singleton:"true"\`` |
| `scope:"prototype"` | Set scope | `_ string \`scope:"prototype"\`` |
| `init-method:"methodName"` | Specify initialization method | `_ string \`init-method:"Connect"\`` |
| `destroy-method:"methodName"` | Specify destruction method | `_ string \`destroy-method:"Close"\`` |

### 🚀 Running Examples

#### Basic Example
```bash
go run examples/basic/main.go
```

#### Web Application Example
```bash
go run examples/web/main.go
```
Then visit http://localhost:8080

#### Logging System Example
```bash
go run examples/logging/main.go
```

#### Run Tests
```bash
go test ./tests/ -v
```

### 📚 Documentation

- [Architecture Design](docs/architecture.md) - Detailed architecture design and implementation principles
- [Usage Guide](docs/usage.md) - Complete usage guide and best practices
- [Logging System](docs/logging.md) - Logging system usage guide and configuration
- [Performance Report](docs/performance.md) - Performance test results and optimization recommendations

### 🤝 Contributing

Issues and Pull Requests are welcome!

### 📄 License

MIT License

---

## Chinese

基于Spring的IoC/DI容器实现思路，使用Go语言的反射、接口和标签实现的轻量级依赖注入框架。

## ✨ 特性

- 🚀 **自动依赖注入** - 基于反射的零配置依赖注入
- 🏷️ **标签驱动** - 使用Go结构体标签进行配置
- 🔄 **生命周期管理** - 完整的Bean初始化和销毁流程
- 📦 **组件扫描** - 自动发现和注册组件
- 🎯 **接口绑定** - 接口和实现的自动绑定
- 🛠️ **多种作用域** - 支持单例和原型模式
- 🧵 **线程安全** - 并发安全的容器实现
- 🔧 **类型安全** - 基于Go强类型系统的安全保证
- 📝 **事件日志** - 完整的框架运行事件记录和监控

## 📁 项目结构

```
gospring/
├── container/          # 核心容器实现
├── context/           # 应用上下文
├── scanner/           # 组件扫描器
├── lifecycle/         # 生命周期管理
├── logging/           # 日志系统
├── annotations/       # 注解和标签处理
├── examples/          # 示例代码
│   ├── basic/        # 基础使用示例
│   ├── web/          # Web应用示例
│   └── logging/      # 日志系统示例
├── tests/            # 单元测试
├── docs/             # 文档
└── README.md
```

## 快速开始

```go
package main

import (
    "fmt"
    "gospring/container"
)

// 定义服务接口
type UserService interface {
    GetUser(id int) string
}

// 实现服务
type UserServiceImpl struct {
    Repository UserRepository `inject:""`
}

func (u *UserServiceImpl) GetUser(id int) string {
    return u.Repository.FindById(id)
}

// 定义仓库接口
type UserRepository interface {
    FindById(id int) string
}

// 实现仓库
type UserRepositoryImpl struct{}

func (u *UserRepositoryImpl) FindById(id int) string {
    return fmt.Sprintf("User-%d", id)
}

func main() {
    // 创建容器
    c := container.NewContainer()
    
    // 注册组件
    c.RegisterSingleton("userRepository", &UserRepositoryImpl{})
    c.RegisterSingleton("userService", &UserServiceImpl{})
    
    // 获取服务
    userService := c.GetBean("userService").(UserService)
    result := userService.GetUser(1)
    fmt.Println(result) // 输出: User-1
}
```

## 核心概念

### 1. 容器 (Container)
负责管理所有组件的生命周期和依赖关系。

### 2. 依赖注入 (Dependency Injection)
通过`inject`标签自动注入依赖。

### 3. 组件扫描 (Component Scanning)
自动发现和注册带有特定标签的组件。

### 4. 生命周期管理
支持组件的初始化和销毁回调。

### 5. 日志系统
提供事件驱动的日志记录，支持多种日志器实现和灵活的日志配置。

## 🏷️ 标签说明

| 标签 | 说明 | 示例 |
|------|------|------|
| `inject:""` | 标记需要注入的字段 | `Repository UserRepo \`inject:""\`` |
| `inject:"beanName"` | 指定注入的Bean名称 | `Cache CacheService \`inject:"redisCache"\`` |
| `component:""` | 标记为组件 | `_ string \`component:"userService"\`` |
| `singleton:"true"` | 标记为单例模式 | `_ string \`singleton:"true"\`` |
| `scope:"prototype"` | 设置作用域 | `_ string \`scope:"prototype"\`` |
| `init-method:"methodName"` | 指定初始化方法 | `_ string \`init-method:"Connect"\`` |
| `destroy-method:"methodName"` | 指定销毁方法 | `_ string \`destroy-method:"Close"\`` |

## 🚀 运行示例

### 基础示例
```bash
go run examples/basic/main.go
```

### Web应用示例
```bash
go run examples/web/main.go
```
然后访问 http://localhost:8080

### 日志系统示例
```bash
go run examples/logging/main.go
```

### 运行测试
```bash
go test ./tests/ -v
```

## 📚 文档

- [架构设计](docs/architecture.md) - 详细的架构设计和实现原理
- [使用指南](docs/usage.md) - 完整的使用指南和最佳实践
- [日志系统](docs/logging.md) - 日志系统的使用指南和配置说明
- [性能报告](docs/performance.md) - 性能测试结果和优化建议

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License