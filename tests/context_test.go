package tests

import (
	"reflect"
	"testing"
	"gospring/context"
)

// 测试用的组件
type TestUserService struct {
	Repository *TestUserRepository `inject:"userRepository"`
	name       string
	
	_ string `component:"userService" singleton:"true"`
}

func (s *TestUserService) GetUser(id int) string {
	return s.Repository.FindUser(id)
}

func (s *TestUserService) Init() error {
	s.name = "initialized"
	return nil
}

func (s *TestUserService) PostConstruct() error {
	return nil
}

type TestUserRepository struct {
	data map[int]string
	
	_ string `component:"userRepository" singleton:"true"`
}

func (r *TestUserRepository) FindUser(id int) string {
	return r.data[id]
}

func (r *TestUserRepository) Init() error {
	r.data = map[int]string{
		1: "用户1",
		2: "用户2",
	}
	return nil
}

func TestApplicationContext_RegisterAndStart(t *testing.T) {
	ctx := context.NewApplicationContext()
	
	// 创建组件
	userRepo := &TestUserRepository{}
	userService := &TestUserService{}
	
	// 注册组件
	err := ctx.RegisterComponents(userRepo, userService)
	if err != nil {
		t.Fatalf("注册组件失败: %v", err)
	}
	
	// 启动上下文
	err = ctx.Start()
	if err != nil {
		t.Fatalf("启动上下文失败: %v", err)
	}
	
	// 验证上下文状态
	if !ctx.IsStarted() {
		t.Error("上下文应该处于启动状态")
	}
	
	// 获取服务并测试
	service := ctx.GetBean("userService")
	if service == nil {
		t.Fatal("无法获取userService")
	}
	
	userSvc, ok := service.(*TestUserService)
	if !ok {
		t.Fatal("服务类型不正确")
	}
	
	// 测试依赖注入是否成功
	repoVal := reflect.ValueOf(userSvc.Repository)
	if !repoVal.IsValid() || (repoVal.Kind() == reflect.Ptr && repoVal.IsNil()) {
		t.Error("依赖注入失败")
	}
	
	// 测试服务功能
	user := userSvc.GetUser(1)
	if user != "用户1" {
		t.Errorf("期望 '用户1', 得到 '%s'", user)
	}
	
	// 测试生命周期
	if userSvc.name != "initialized" {
		t.Error("初始化方法未被调用")
	}
	
	// 停止上下文
	err = ctx.Stop()
	if err != nil {
		t.Fatalf("停止上下文失败: %v", err)
	}
	
	if ctx.IsStarted() {
		t.Error("上下文应该处于停止状态")
	}
}

func TestApplicationContext_GetBeanByType(t *testing.T) {
	ctx := context.NewApplicationContext()
	
	userService := &TestUserService{}
	ctx.RegisterComponent(userService)
	
	ctx.Start()
	
	// 通过类型获取
	serviceType := reflect.TypeOf(userService)
	bean := ctx.GetBeanByType(serviceType)
	if bean == nil {
		t.Fatal("无法通过类型获取Bean")
	}
	
	svc, ok := bean.(*TestUserService)
	if !ok {
		t.Fatal("Bean类型不正确")
	}
	
	if svc != userService {
		t.Error("返回的Bean实例不正确")
	}
}

func TestApplicationContext_ListBeans(t *testing.T) {
	ctx := context.NewApplicationContext()
	
	userRepo := &TestUserRepository{}
	userService := &TestUserService{}
	
	ctx.RegisterComponents(userRepo, userService)
	ctx.Start()
	
	beans := ctx.ListBeans()
	if len(beans) != 2 {
		t.Errorf("期望2个Bean, 得到%d个", len(beans))
	}
	
	expectedBeans := []string{"userRepository", "userService"}
	for _, expected := range expectedBeans {
		found := false
		for _, bean := range beans {
			if bean == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("未找到期望的Bean: %s", expected)
		}
	}
}

func TestApplicationContext_HasBean(t *testing.T) {
	ctx := context.NewApplicationContext()
	
	userService := &TestUserService{}
	ctx.RegisterComponent(userService)
	ctx.Start()
	
	if !ctx.HasBean("userService") {
		t.Error("应该存在userService Bean")
	}
	
	if ctx.HasBean("nonExistentBean") {
		t.Error("不应该存在nonExistentBean")
	}
}

func TestApplicationContext_GetBeansOfType(t *testing.T) {
	ctx := context.NewApplicationContext()
	
	userService1 := &TestUserService{}
	userService2 := &TestUserService{}
	
	ctx.RegisterBean("userService1", userService1)
	ctx.RegisterBean("userService2", userService2)
	ctx.Start()
	
	serviceType := reflect.TypeOf(userService1)
	beans := ctx.GetBeansOfType(serviceType)
	
	if len(beans) != 2 {
		t.Errorf("期望2个相同类型的Bean, 得到%d个", len(beans))
	}
	
	if _, exists := beans["userService1"]; !exists {
		t.Error("应该包含userService1")
	}
	
	if _, exists := beans["userService2"]; !exists {
		t.Error("应该包含userService2")
	}
}

func TestApplicationContext_AutoWire(t *testing.T) {
	ctx := context.NewApplicationContext()
	
	// 注册依赖
	userRepo := &TestUserRepository{}
	ctx.RegisterComponent(userRepo)
	ctx.Start()
	
	// 创建需要自动装配的对象
	userService := &TestUserService{}
	
	// 执行自动装配
	err := ctx.AutoWire(userService)
	if err != nil {
		t.Fatalf("自动装配失败: %v", err)
	}
	
	// 验证装配结果
	repoVal := reflect.ValueOf(userService.Repository)
	if !repoVal.IsValid() || (repoVal.Kind() == reflect.Ptr && repoVal.IsNil()) {
		t.Error("自动装配失败，Repository为nil")
	}
}

func TestApplicationContext_Refresh(t *testing.T) {
	ctx := context.NewApplicationContext()
	
	userService := &TestUserService{}
	ctx.RegisterComponent(userService)
	
	// 启动
	err := ctx.Start()
	if err != nil {
		t.Fatalf("启动失败: %v", err)
	}
	
	if !ctx.IsStarted() {
		t.Error("上下文应该处于启动状态")
	}
	
	// 刷新
	err = ctx.Refresh()
	if err != nil {
		t.Fatalf("刷新失败: %v", err)
	}
	
	if !ctx.IsStarted() {
		t.Error("刷新后上下文应该处于启动状态")
	}
}