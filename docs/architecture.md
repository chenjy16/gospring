# GoSpring 架构设计

## 整体架构

GoSpring 框架采用分层架构设计，主要包含以下核心模块：

```
┌─────────────────────────────────────────┐
│              Application Layer           │
│         (examples/basic, web)           │
├─────────────────────────────────────────┤
│              Context Layer              │
│            (context package)           │
├─────────────────────────────────────────┤
│              Service Layer              │
│  ┌─────────────┬─────────────┬─────────┐ │
│  │   Scanner   │ Lifecycle   │ Annotations│ │
│  │   Package   │  Package    │ Package │ │
│  └─────────────┴─────────────┴─────────┘ │
├─────────────────────────────────────────┤
│              Core Layer                 │
│           (container package)          │
└─────────────────────────────────────────┘
```

## 核心组件

### 1. Container (容器)
- **职责**: Bean的注册、管理和获取
- **特性**: 
  - 支持单例和原型模式
  - 线程安全的Bean管理
  - 类型映射和接口绑定
  - 依赖注入执行

### 2. ApplicationContext (应用上下文)
- **职责**: 整合各个模块，提供统一的API
- **特性**:
  - 生命周期管理
  - 组件扫描集成
  - 启动/停止控制
  - 泛型Bean获取

### 3. ComponentScanner (组件扫描器)
- **职责**: 自动发现和注册组件
- **特性**:
  - 基于标签的组件识别
  - 批量组件注册
  - 接口实现自动绑定

### 4. LifecycleManager (生命周期管理器)
- **职责**: Bean的初始化和销毁管理
- **特性**:
  - 多种初始化回调支持
  - 销毁顺序管理
  - 错误处理和恢复

### 5. Annotations (注解系统)
- **职责**: 标签解析和组件识别
- **特性**:
  - 多种组件类型支持
  - 灵活的标签配置
  - 反射工具集成

## 设计模式

### 1. 依赖注入 (Dependency Injection)
```go
type UserService struct {
    Repository UserRepository `inject:"userRepository"`
    Cache      CacheService   `inject:"cacheService"`
}
```

### 2. 工厂模式 (Factory Pattern)
```go
// 容器作为Bean工厂
func (c *Container) GetBean(name string) interface{}
func (c *Container) createNewInstance(beanDef *BeanDefinition) interface{}
```

### 3. 单例模式 (Singleton Pattern)
```go
type BeanDefinition struct {
    Singleton bool
    Instance  interface{}
    mutex     sync.RWMutex
}
```

### 4. 观察者模式 (Observer Pattern)
```go
// 生命周期回调
type Initializer interface {
    Init() error
}

type PostConstruct interface {
    PostConstruct() error
}
```

## 反射机制

### 1. 类型检查
```go
func (c *Container) registerBean(name string, instance interface{}, singleton bool) error {
    val := reflect.ValueOf(instance)
    typ := reflect.TypeOf(instance)
    
    if typ.Kind() == reflect.Ptr {
        typ = typ.Elem()
    }
    // ...
}
```

### 2. 字段注入
```go
func (c *Container) InjectDependencies(instance interface{}) error {
    val := reflect.ValueOf(instance)
    typ := val.Type()
    
    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        fieldType := typ.Field(i)
        
        injectTag := fieldType.Tag.Get("inject")
        if injectTag != "" && field.CanSet() {
            // 执行注入
        }
    }
}
```

### 3. 方法调用
```go
func (lm *LifecycleManager) callInitMethod(instance interface{}) error {
    val := reflect.ValueOf(instance)
    
    method := val.MethodByName("Init")
    if method.IsValid() {
        results := method.Call(nil)
        // 处理返回值
    }
}
```

## 标签系统

### 1. 组件标签
```go
type UserService struct {
    _ string `component:"userService" singleton:"true"`
}
```

### 2. 注入标签
```go
type UserController struct {
    Service UserService `inject:"userService"`
}
```

### 3. 生命周期标签
```go
type DatabaseService struct {
    _ string `init-method:"Connect" destroy-method:"Disconnect"`
}
```

## 线程安全

### 1. 容器级别
```go
type Container struct {
    beans map[string]*BeanDefinition
    mutex sync.RWMutex
}
```

### 2. Bean级别
```go
type BeanDefinition struct {
    Instance interface{}
    mutex    sync.RWMutex
}
```

## 扩展点

### 1. 自定义组件类型
```go
type CustomComponent interface {
    Component
    CustomMethod() error
}
```

### 2. 自定义生命周期回调
```go
type CustomLifecycle interface {
    BeforeInit() error
    AfterInit() error
}
```

### 3. 自定义注入逻辑
```go
type CustomInjector interface {
    Inject(field reflect.StructField, value reflect.Value) error
}
```

## 性能考虑

### 1. 反射缓存
- 类型信息缓存
- 方法查找缓存
- 字段映射缓存

### 2. 延迟初始化
- 按需创建Bean
- 延迟依赖解析
- 懒加载支持

### 3. 内存管理
- 弱引用支持
- 资源清理
- 垃圾回收友好