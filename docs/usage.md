# GoSpring 使用指南

## 快速开始

### 1. 基本概念

GoSpring 是一个基于 Spring IoC/DI 思想的 Go 语言框架，主要特性包括：

- **依赖注入**: 通过 `inject` 标签自动注入依赖
- **组件管理**: 通过 `component` 标签标记组件
- **生命周期**: 支持初始化和销毁回调
- **类型安全**: 基于 Go 的强类型系统

### 2. 基本用法

#### 定义组件

```go
// 定义服务接口
type UserService interface {
    GetUser(id int) *User
}

// 实现服务
type UserServiceImpl struct {
    Repository UserRepository `inject:"userRepository"`
    
    // 组件标记
    _ string `component:"userService" singleton:"true"`
}

func (s *UserServiceImpl) GetUser(id int) *User {
    return s.Repository.FindById(id)
}

// 生命周期方法
func (s *UserServiceImpl) Init() error {
    fmt.Println("UserService 初始化")
    return nil
}
```

#### 创建和使用容器

```go
func main() {
    // 创建应用上下文
    ctx := context.NewApplicationContext()
    
    // 注册组件
    userRepo := &UserRepositoryImpl{}
    userService := &UserServiceImpl{}
    
    ctx.RegisterComponents(userRepo, userService)
    
    // 启动上下文
    ctx.Start()
    
    // 获取并使用服务
    service := ctx.GetBean("userService").(UserService)
    user := service.GetUser(1)
    
    // 停止上下文
    ctx.Stop()
}
```

## 详细功能

### 1. 依赖注入

#### 基本注入
```go
type OrderService struct {
    UserService    UserService    `inject:"userService"`
    ProductService ProductService `inject:"productService"`
}
```

#### 按类型注入
```go
type OrderService struct {
    UserService UserService `inject:""` // 自动按类型查找
}
```

#### 可选注入
```go
type OrderService struct {
    CacheService CacheService `inject:"cacheService,optional"`
}
```

### 2. 组件注册

#### 手动注册
```go
ctx := context.NewApplicationContext()

// 注册单例
ctx.RegisterBean("userService", &UserServiceImpl{})

// 注册原型
container := ctx.GetContainer()
container.RegisterPrototype("prototypeBean", &PrototypeBean{})
```

#### 自动扫描
```go
// 使用组件标签
type UserService struct {
    _ string `component:"userService" singleton:"true"`
}

// 批量注册
ctx.RegisterComponents(
    &UserServiceImpl{},
    &UserRepositoryImpl{},
    &CacheServiceImpl{},
)
```

#### 接口绑定
```go
// 绑定接口和实现
userServiceType := reflect.TypeOf((*UserService)(nil)).Elem()
ctx.RegisterByInterface(userServiceType, &UserServiceImpl{}, "userService")

// 通过接口获取
service := ctx.GetBeanByType(userServiceType).(UserService)
```

### 3. 生命周期管理

#### 初始化回调
```go
type DatabaseService struct {
    conn *sql.DB
}

// 方法1: 实现 Initializer 接口
func (s *DatabaseService) Init() error {
    var err error
    s.conn, err = sql.Open("mysql", "dsn")
    return err
}

// 方法2: 实现 PostConstruct 接口
func (s *DatabaseService) PostConstruct() error {
    return s.setupTables()
}

// 方法3: 使用标签指定方法
type DatabaseService struct {
    _ string `init-method:"Connect"`
}

func (s *DatabaseService) Connect() error {
    // 连接逻辑
    return nil
}
```

#### 销毁回调
```go
type DatabaseService struct {
    conn *sql.DB
}

// 方法1: 实现 Destroyer 接口
func (s *DatabaseService) Destroy() error {
    if s.conn != nil {
        return s.conn.Close()
    }
    return nil
}

// 方法2: 实现 PreDestroy 接口
func (s *DatabaseService) PreDestroy() error {
    return s.cleanup()
}

// 方法3: 使用标签指定方法
type DatabaseService struct {
    _ string `destroy-method:"Disconnect"`
}
```

#### Bean名称感知
```go
type LoggingService struct {
    beanName string
}

func (s *LoggingService) SetBeanName(name string) {
    s.beanName = name
    fmt.Printf("Bean名称设置为: %s\n", name)
}
```

### 4. 作用域管理

#### 单例模式（默认）
```go
type SingletonService struct {
    _ string `component:"singletonService" singleton:"true"`
}
```

#### 原型模式
```go
type PrototypeService struct {
    _ string `component:"prototypeService" singleton:"false"`
}

// 或者使用 scope 标签
type PrototypeService struct {
    _ string `component:"prototypeService" scope:"prototype"`
}
```

### 5. 高级用法

#### 泛型Bean获取（Go 1.18+）
```go
// 类型安全的Bean获取
userService := context.GetBeanT[UserService](ctx, "userService")
```

#### 获取同类型的所有Bean
```go
userServiceType := reflect.TypeOf((*UserService)(nil)).Elem()
services := ctx.GetBeansOfType(userServiceType)

for name, service := range services {
    fmt.Printf("服务: %s, 类型: %T\n", name, service)
}
```

#### 自动装配外部对象
```go
// 对已存在的对象执行依赖注入
externalObject := &ExternalService{}
ctx.AutoWire(externalObject)
```

#### 动态Bean创建
```go
// 使用工厂函数创建Bean
ctx.CreateBean("dynamicBean", func() interface{} {
    return &DynamicService{
        timestamp: time.Now(),
    }
})
```

## Web应用集成

### 1. HTTP控制器
```go
type ProductController struct {
    ProductService ProductService `inject:"productService"`
    
    _ string `component:"productController"`
}

func (c *ProductController) SetupRoutes() {
    http.HandleFunc("/products", c.handleProducts)
    http.HandleFunc("/products/", c.handleProduct)
}

func (c *ProductController) handleProducts(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        products := c.ProductService.GetAllProducts()
        json.NewEncoder(w).Encode(products)
    case "POST":
        // 处理创建请求
    }
}
```

### 2. 中间件集成
```go
type AuthMiddleware struct {
    UserService UserService `inject:"userService"`
    
    _ string `component:"authMiddleware"`
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 认证逻辑
        if m.UserService.ValidateToken(r.Header.Get("Authorization")) {
            next.ServeHTTP(w, r)
        } else {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
        }
    })
}
```

## 最佳实践

### 1. 组件设计
- 保持组件的单一职责
- 优先使用接口而不是具体实现
- 避免循环依赖

### 2. 依赖注入
- 使用构造函数注入（通过Init方法）
- 避免在字段注入中使用复杂逻辑
- 明确标记可选依赖

### 3. 生命周期管理
- 在Init方法中进行资源初始化
- 在Destroy方法中进行资源清理
- 处理初始化和销毁中的错误

### 4. 性能优化
- 合理使用单例和原型模式
- 避免过度使用反射
- 缓存频繁访问的Bean

### 5. 错误处理
```go
func main() {
    ctx := context.NewApplicationContext()
    
    // 注册组件时处理错误
    if err := ctx.RegisterComponents(services...); err != nil {
        log.Fatalf("注册组件失败: %v", err)
    }
    
    // 启动时处理错误
    if err := ctx.Start(); err != nil {
        log.Fatalf("启动上下文失败: %v", err)
    }
    
    // 确保资源清理
    defer func() {
        if err := ctx.Stop(); err != nil {
            log.Printf("停止上下文时出错: %v", err)
        }
    }()
}
```

## 常见问题

### 1. 循环依赖
```go
// 问题：A依赖B，B依赖A
type ServiceA struct {
    ServiceB ServiceB `inject:"serviceB"`
}

type ServiceB struct {
    ServiceA ServiceA `inject:"serviceA"`
}

// 解决方案：引入第三个服务或使用事件机制
```

### 2. 接口注入失败
```go
// 确保接口类型正确注册
userServiceType := reflect.TypeOf((*UserService)(nil)).Elem()
ctx.RegisterByInterface(userServiceType, implementation, "beanName")
```

### 3. 生命周期方法未调用
```go
// 确保实现了正确的接口
type MyService struct{}

func (s *MyService) Init() error {  // 必须返回error
    // 初始化逻辑
    return nil
}
```