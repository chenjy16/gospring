package tests

import (
	"reflect"
	"testing"
	"time"
	"gospring/container"
	"gospring/logging"
	"gospring/scanner"
	"github.com/stretchr/testify/assert"
)

// 测试用的组件结构体
type ScanTestService struct {
	Name string `component:"testService"`
}

func (s *ScanTestService) GetName() string {
	return s.Name
}

type ScanTestRepository struct {
	Data string `component:"testRepository" singleton:"true"`
}

type ScanTestController struct {
	Service *ScanTestService `inject:"testService"`
	Info    string           `component:"testController" singleton:"false"`
}

// 为不同测试使用的组件结构体，避免名称冲突
type ScanTestService2 struct {
	Name string `component:"testService2"`
}

type ScanTestService3 struct {
	Name string `component:"testService3"`
}

type ScanTestService4 struct {
	Name string `component:"testService4"`
}

type ScanPlainStruct struct {
	Name string
}

type ScanTestComponent struct {
	Value string `component:"true"`
}

type ScanCustomNameComponent struct {
	Value string `component:"customName"`
}

// TestNewComponentScanner 测试创建组件扫描器
func TestNewComponentScanner(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	assert.NotNil(t, s)
	assert.NotNil(t, s.GetLogger())
}

// TestNewComponentScannerWithLogger 测试创建带日志器的组件扫描器
func TestNewComponentScannerWithLogger(t *testing.T) {
	c := container.NewContainer()
	logger := logging.NopLogger
	s := scanner.NewComponentScannerWithLogger(c, logger)
	
	assert.NotNil(t, s)
	assert.Equal(t, logger, s.GetLogger())
}

// TestSetLogger 测试设置日志器
func TestSetLogger(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	logger := logging.NopLogger
	
	s.SetLogger(logger)
	assert.Equal(t, logger, s.GetLogger())
}

// TestAddPackage 测试添加包
func TestAddPackage(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	s.AddPackage("github.com/test/package1")
	s.AddPackage("github.com/test/package2")
	
	// 由于 packages 字段是私有的，我们无法直接测试，但可以确保不会panic
	assert.NotNil(t, s)
}

// TestScanComponent_Success 测试成功扫描组件
func TestScanComponent_Success(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	service := &ScanTestService{Name: "test"}
	err := s.ScanComponent(service)
	
	assert.NoError(t, err)
	
	// 检查组件是否被注册
	bean := c.GetBean("testService")
	assert.NotNil(t, bean)
}

// TestScanComponent_NonComponent 测试扫描非组件
func TestScanComponent_NonComponent(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	nonComponent := &ScanPlainStruct{Name: "test"}
	err := s.ScanComponent(nonComponent)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not a component")
}

// TestScanComponent_CustomName 测试自定义组件名称
func TestScanComponent_CustomName(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	component := &ScanCustomNameComponent{Value: "test"}
	err := s.ScanComponent(component)
	
	assert.NoError(t, err)
	
	// 检查组件是否以自定义名称注册
	bean := c.GetBean("customName")
	assert.NotNil(t, bean)
}

// TestScanComponent_Singleton 测试单例组件
func TestScanComponent_Singleton(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	repository := &ScanTestRepository{Data: "test"}
	err := s.ScanComponent(repository)
	
	assert.NoError(t, err)
	
	// 检查组件是否被注册为单例
	bean1 := c.GetBean("testRepository")
	assert.NotNil(t, bean1)
	bean2 := c.GetBean("testRepository")
	assert.NotNil(t, bean2)
	assert.Same(t, bean1, bean2) // 应该是同一个实例
}

// TestScanComponent_Prototype 测试原型组件
func TestScanComponent_Prototype(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	controller := &ScanTestController{Info: "test"}
	err := s.ScanComponent(controller)
	
	assert.NoError(t, err)
	
	// 检查组件是否被注册
	bean := c.GetBean("testController")
	assert.NotNil(t, bean)
}

// TestScanAndRegister 测试批量扫描和注册
func TestScanAndRegister(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	service := &ScanTestService2{Name: "service"}
	repository := &ScanTestRepository{Data: "repository"}
	
	err := s.ScanAndRegister(service, repository)
	assert.NoError(t, err)
	
	// 检查两个组件都被注册
	bean1 := c.GetBean("testService2")
	assert.NotNil(t, bean1)
	
	bean2 := c.GetBean("testRepository")
	assert.NotNil(t, bean2)
}

// TestScanAndRegister_WithError 测试批量扫描时的错误处理
func TestScanAndRegister_WithError(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	service := &ScanTestService3{Name: "service"}
	nonComponent := &ScanPlainStruct{Name: "non"}
	
	err := s.ScanAndRegister(service, nonComponent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to scan component")
	
	// 第一个组件应该成功注册
	bean := c.GetBean("testService3")
	assert.NotNil(t, bean)
}

// TestAutoScan 测试自动扫描
func TestAutoScan(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	// 创建一个包含组件字段的结构体
	type ScanContainer struct {
		Service    *ScanTestService4   `component:"true"`
		Repository *ScanTestRepository `component:"true"`
		NonComp    *ScanPlainStruct   // 没有component标签
	}
	
	scanContainer := &ScanContainer{
		Service:    &ScanTestService4{Name: "auto"},
		Repository: &ScanTestRepository{Data: "auto"},
		NonComp:    &ScanPlainStruct{Name: "non"},
	}
	
	err := s.AutoScan(scanContainer)
	assert.NoError(t, err)
	
	// 检查组件是否被注册
	bean1 := c.GetBean("testService4")
	assert.NotNil(t, bean1)
	
	bean2 := c.GetBean("testRepository")
	assert.NotNil(t, bean2)
	
	// NonComponent 不应该被注册
	bean3 := c.GetBean("noncomponent")
	assert.Nil(t, bean3)
}

// TestRegisterWithInterface 测试按接口注册
func TestRegisterWithInterface(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	// 定义一个接口
	type ScanServiceInterface interface {
		GetName() string
	}
	
	// 让 ScanTestService 实现接口
	service := &ScanTestService{Name: "interface_test"}
	
	interfaceType := reflect.TypeOf((*ScanServiceInterface)(nil)).Elem()
	err := s.RegisterWithInterface(interfaceType, service, "serviceInterface")
	
	assert.NoError(t, err)
	
	// 检查是否可以按接口获取
	bean := c.GetBeanByType(interfaceType)
	assert.NotNil(t, bean)
}

// TestScanPackageComponents 测试包扫描
func TestScanPackageComponents(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	// 这个方法目前只是打印信息，不会出错
	err := s.ScanPackageComponents()
	assert.NoError(t, err)
}

// TestComponentNamingConventions 测试组件命名约定
func TestComponentNamingConventions(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	// 测试不同的命名约定
	type UserService struct {
		Name string
	}
	
	type OrderRepository struct {
		Data string
	}
	
	type ProductController struct {
		Info string
	}
	
	type CustomComponent struct {
		Value string
	}
	
	userService := &UserService{Name: "user"}
	orderRepository := &OrderRepository{Data: "order"}
	productController := &ProductController{Info: "product"}
	customComponent := &CustomComponent{Value: "custom"}
	
	// 这些应该根据命名约定被识别为组件
	err := s.ScanComponent(userService)
	assert.NoError(t, err)
	
	err = s.ScanComponent(orderRepository)
	assert.NoError(t, err)
	
	err = s.ScanComponent(productController)
	assert.NoError(t, err)
	
	err = s.ScanComponent(customComponent)
	assert.NoError(t, err)
	
	// 检查组件是否被正确注册
	bean1 := c.GetBean("userservice")
	assert.NotNil(t, bean1)
	
	bean2 := c.GetBean("orderrepository")
	assert.NotNil(t, bean2)
	
	bean3 := c.GetBean("productcontroller")
	assert.NotNil(t, bean3)
	
	bean4 := c.GetBean("customcomponent")
	assert.NotNil(t, bean4)
}

// TestScanComponent_PointerHandling 测试指针处理
func TestScanComponent_PointerHandling(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	// 测试传入指针
	service := &ScanTestService{Name: "pointer"}
	err := s.ScanComponent(service)
	assert.NoError(t, err)
	
	// 测试传入值（会被转换为指针）
	service2 := ScanTestService{Name: "value"}
	err = s.ScanComponent(service2)
	assert.Error(t, err) // 第二次注册同名组件应该失败
	assert.Contains(t, err.Error(), "already exists")
	
	// 第一次注册的组件应该存在
	bean := c.GetBean("testService")
	assert.NotNil(t, bean)
}

// TestEventTiming 测试事件时间记录
func TestEventTiming(t *testing.T) {
	c := container.NewContainer()
	s := scanner.NewComponentScanner(c)
	
	service := &ScanTestService{Name: "timing"}
	
	start := time.Now()
	err := s.ScanComponent(service)
	end := time.Now()
	
	assert.NoError(t, err)
	
	// 检查扫描时间是否合理
	duration := end.Sub(start)
	assert.True(t, duration > 0)
	assert.True(t, duration < time.Second) // 应该很快完成
}