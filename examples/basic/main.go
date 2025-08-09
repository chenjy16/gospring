package main

import (
	"fmt"
	"log"
	"reflect"
	"gospring/context"
)

// 定义服务接口
type UserService interface {
	GetUser(id int) *User
	CreateUser(name, email string) *User
}

// 定义仓库接口
type UserRepository interface {
	FindById(id int) *User
	Save(user *User) error
}

// 定义缓存接口
type CacheService interface {
	Get(key string) interface{}
	Set(key string, value interface{})
}

// User 用户模型
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	Repository UserRepository `inject:"userRepository"`
	Cache      CacheService   `inject:"cacheService"`
	
	// 组件标记
	_ string `component:"userService" singleton:"true"`
}

// 实现UserService接口
func (s *UserServiceImpl) GetUser(id int) *User {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("user:%d", id)
	if cached := s.Cache.Get(cacheKey); cached != nil {
		if user, ok := cached.(*User); ok {
			fmt.Printf("从缓存获取用户: %+v\n", user)
			return user
		}
	}

	// 从仓库获取
	user := s.Repository.FindById(id)
	if user != nil {
		s.Cache.Set(cacheKey, user)
		fmt.Printf("从数据库获取用户: %+v\n", user)
	}
	
	return user
}

func (s *UserServiceImpl) CreateUser(name, email string) *User {
	user := &User{
		ID:    len(name), // 简单的ID生成
		Name:  name,
		Email: email,
	}
	
	s.Repository.Save(user)
	
	// 缓存新用户
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	s.Cache.Set(cacheKey, user)
	
	return user
}

// 实现生命周期接口
func (s *UserServiceImpl) Init() error {
	fmt.Println("UserService 初始化完成")
	return nil
}

func (s *UserServiceImpl) PostConstruct() error {
	fmt.Println("UserService 构造后处理完成")
	return nil
}

// UserRepositoryImpl 用户仓库实现
type UserRepositoryImpl struct {
	users map[int]*User
	
	// 组件标记
	_ string `component:"userRepository" singleton:"true"`
}

func (r *UserRepositoryImpl) FindById(id int) *User {
	return r.users[id]
}

func (r *UserRepositoryImpl) Save(user *User) error {
	r.users[user.ID] = user
	fmt.Printf("保存用户到数据库: %+v\n", user)
	return nil
}

// 初始化方法
func (r *UserRepositoryImpl) Init() error {
	r.users = make(map[int]*User)
	// 初始化一些测试数据
	r.users[1] = &User{ID: 1, Name: "张三", Email: "zhangsan@example.com"}
	r.users[2] = &User{ID: 2, Name: "李四", Email: "lisi@example.com"}
	fmt.Println("UserRepository 初始化完成，加载测试数据")
	return nil
}

// CacheServiceImpl 缓存服务实现
type CacheServiceImpl struct {
	cache map[string]interface{}
	
	// 组件标记
	_ string `component:"cacheService" singleton:"true"`
}

func (c *CacheServiceImpl) Get(key string) interface{} {
	return c.cache[key]
}

func (c *CacheServiceImpl) Set(key string, value interface{}) {
	c.cache[key] = value
}

func (c *CacheServiceImpl) Init() error {
	c.cache = make(map[string]interface{})
	fmt.Println("CacheService 初始化完成")
	return nil
}

func (c *CacheServiceImpl) Destroy() error {
	c.cache = nil
	fmt.Println("CacheService 销毁完成")
	return nil
}

func main() {
	fmt.Println("=== GoSpring 框架演示 ===")
	
	// 创建应用上下文
	ctx := context.NewApplicationContext()
	
	// 创建组件实例
	userRepo := &UserRepositoryImpl{}
	cacheService := &CacheServiceImpl{}
	userService := &UserServiceImpl{}
	
	// 注册组件
	fmt.Println("\n1. 注册组件...")
	if err := ctx.RegisterComponents(userRepo, cacheService, userService); err != nil {
		log.Fatalf("注册组件失败: %v", err)
	}
	
	// 也可以通过接口注册
	fmt.Println("2. 通过接口注册组件...")
	userServiceType := reflect.TypeOf((*UserService)(nil)).Elem()
	if err := ctx.RegisterByInterface(userServiceType, userService, "userServiceInterface"); err != nil {
		log.Fatalf("通过接口注册失败: %v", err)
	}
	
	// 启动上下文
	fmt.Println("\n3. 启动应用上下文...")
	if err := ctx.Start(); err != nil {
		log.Fatalf("启动上下文失败: %v", err)
	}
	
	// 获取服务并使用
	fmt.Println("\n4. 使用服务...")
	
	// 方式1：通过名称获取
	service := ctx.GetBean("userService").(UserService)
	
	// 测试获取用户
	fmt.Println("\n--- 测试获取用户 ---")
	user1 := service.GetUser(1)
	fmt.Printf("获取到用户1: %+v\n", user1)
	
	// 再次获取同一用户（应该从缓存获取）
	user1Again := service.GetUser(1)
	fmt.Printf("再次获取用户1: %+v\n", user1Again)
	
	// 测试创建用户
	fmt.Println("\n--- 测试创建用户 ---")
	newUser := service.CreateUser("王五", "wangwu@example.com")
	fmt.Printf("创建新用户: %+v\n", newUser)
	
	// 获取新创建的用户
	retrievedUser := service.GetUser(newUser.ID)
	fmt.Printf("获取新创建的用户: %+v\n", retrievedUser)
	
	// 方式2：通过类型获取
	fmt.Println("\n5. 通过类型获取服务...")
	serviceByType := ctx.GetBeanByType(userServiceType).(UserService)
	user2 := serviceByType.GetUser(2)
	fmt.Printf("通过类型获取的服务，获取用户2: %+v\n", user2)
	
	// 显示所有注册的Bean
	fmt.Println("\n6. 显示所有注册的Bean:")
	beans := ctx.ListBeans()
	for _, beanName := range beans {
		beanDef := ctx.GetBeanDefinition(beanName)
		fmt.Printf("  - %s (类型: %v, 单例: %v)\n", 
			beanName, beanDef.Type, beanDef.Singleton)
	}
	
	// 获取指定类型的所有Bean
	fmt.Println("\n7. 获取UserService类型的所有Bean:")
	userServices := ctx.GetBeansOfType(userServiceType)
	for name, bean := range userServices {
		fmt.Printf("  - %s: %T\n", name, bean)
	}
	
	// 停止上下文
	fmt.Println("\n8. 停止应用上下文...")
	if err := ctx.Stop(); err != nil {
		log.Fatalf("停止上下文失败: %v", err)
	}
	
	fmt.Println("\n=== 演示完成 ===")
}