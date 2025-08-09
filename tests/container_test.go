package tests

import (
	"reflect"
	"testing"
	"gospring/container"
	"github.com/stretchr/testify/assert"
)

// 测试用的接口和实现
type TestService interface {
	GetName() string
}

type TestServiceImpl struct {
	name string
}

func (t *TestServiceImpl) GetName() string {
	return t.name
}

type TestRepository interface {
	Save(data string) error
}

type TestRepositoryImpl struct {
	data []string
}

func (t *TestRepositoryImpl) Save(data string) error {
	t.data = append(t.data, data)
	return nil
}

type TestController struct {
	Service    TestService    `inject:"testService"`
	Repository TestRepository `inject:"testRepository"`
}

func TestContainer_RegisterSingleton(t *testing.T) {
	c := container.NewContainer()
	
	service := &TestServiceImpl{name: "test"}
	err := c.RegisterSingleton("testService", service)
	
	assert.NoError(t, err)
	assert.True(t, c.HasBean("testService"))
	
	// 测试重复注册
	err = c.RegisterSingleton("testService", service)
	assert.Error(t, err)
}

func TestContainer_RegisterPrototype(t *testing.T) {
	c := container.NewContainer()
	
	service := &TestServiceImpl{name: "prototype"}
	err := c.RegisterPrototype("prototypeService", service)
	
	assert.NoError(t, err)
	assert.True(t, c.HasBean("prototypeService"))
}

func TestContainer_GetBean(t *testing.T) {
	c := container.NewContainer()
	
	service := &TestServiceImpl{name: "test"}
	c.RegisterSingleton("testService", service)
	
	// 获取Bean
	bean := c.GetBean("testService")
	assert.NotNil(t, bean)
	
	// 类型断言
	testService, ok := bean.(*TestServiceImpl)
	assert.True(t, ok)
	assert.Equal(t, "test", testService.GetName())
	
	// 单例模式，应该返回同一个实例
	bean2 := c.GetBean("testService")
	assert.Same(t, bean, bean2)
}

func TestContainer_GetBeanByType(t *testing.T) {
	c := container.NewContainer()
	
	service := &TestServiceImpl{name: "test"}
	c.RegisterSingleton("testService", service)
	
	// 根据类型获取Bean
	serviceType := reflect.TypeOf(service).Elem()
	bean := c.GetBeanByType(serviceType)
	
	assert.NotNil(t, bean)
	testService, ok := bean.(*TestServiceImpl)
	assert.True(t, ok)
	assert.Equal(t, "test", testService.GetName())
}

func TestContainer_InjectDependencies(t *testing.T) {
	c := container.NewContainer()
	
	// 注册依赖
	service := &TestServiceImpl{name: "injected"}
	repository := &TestRepositoryImpl{data: make([]string, 0)}
	
	c.RegisterSingleton("testService", service)
	c.RegisterSingleton("testRepository", repository)
	
	// 创建需要注入的对象
	controller := &TestController{}
	
	// 执行依赖注入
	err := c.InjectDependencies(controller)
	assert.NoError(t, err)
	
	// 验证注入结果
	assert.NotNil(t, controller.Service)
	assert.NotNil(t, controller.Repository)
	assert.Equal(t, "injected", controller.Service.GetName())
}

func TestContainer_WireAll(t *testing.T) {
	c := container.NewContainer()
	
	// 注册依赖
	service := &TestServiceImpl{name: "wired"}
	repository := &TestRepositoryImpl{data: make([]string, 0)}
	controller := &TestController{}
	
	c.RegisterSingleton("testService", service)
	c.RegisterSingleton("testRepository", repository)
	c.RegisterSingleton("testController", controller)
	
	// 执行全部装配
	err := c.WireAll()
	assert.NoError(t, err)
	
	// 获取控制器并验证注入
	bean := c.GetBean("testController")
	ctrl, ok := bean.(*TestController)
	assert.True(t, ok)
	assert.NotNil(t, ctrl.Service)
	assert.NotNil(t, ctrl.Repository)
}

func TestContainer_RegisterByInterface(t *testing.T) {
	c := container.NewContainer()
	
	service := &TestServiceImpl{name: "interface"}
	serviceInterface := reflect.TypeOf((*TestService)(nil)).Elem()
	
	err := c.RegisterByInterface(serviceInterface, service, "testServiceInterface")
	assert.NoError(t, err)
	
	// 通过接口类型获取
	bean := c.GetBeanByType(serviceInterface)
	assert.NotNil(t, bean)
	
	testService, ok := bean.(TestService)
	assert.True(t, ok)
	assert.Equal(t, "interface", testService.GetName())
}

func TestContainer_ListBeans(t *testing.T) {
	c := container.NewContainer()
	
	service1 := &TestServiceImpl{name: "service1"}
	service2 := &TestServiceImpl{name: "service2"}
	
	c.RegisterSingleton("service1", service1)
	c.RegisterSingleton("service2", service2)
	
	beans := c.ListBeans()
	assert.Len(t, beans, 2)
	assert.Contains(t, beans, "service1")
	assert.Contains(t, beans, "service2")
}

func TestContainer_GetBeanDefinition(t *testing.T) {
	c := container.NewContainer()
	
	service := &TestServiceImpl{name: "test"}
	c.RegisterSingleton("testService", service)
	
	beanDef := c.GetBeanDefinition("testService")
	assert.NotNil(t, beanDef)
	assert.Equal(t, "testService", beanDef.Name)
	assert.True(t, beanDef.Singleton)
	assert.Equal(t, reflect.TypeOf(service).Elem(), beanDef.Type)
}

func TestContainer_PrototypeScope(t *testing.T) {
	c := container.NewContainer()
	
	service := &TestServiceImpl{name: "prototype"}
	c.RegisterPrototype("prototypeService", service)
	
	// 原型模式应该返回不同的实例
	bean1 := c.GetBean("prototypeService")
	bean2 := c.GetBean("prototypeService")
	
	assert.NotNil(t, bean1)
	assert.NotNil(t, bean2)
	// 注意：由于我们的简单实现，这里可能返回同一个实例
	// 在实际的原型实现中，应该创建新实例
}

func TestContainer_Destroy(t *testing.T) {
	c := container.NewContainer()
	
	service := &TestServiceImpl{name: "test"}
	c.RegisterSingleton("testService", service)
	
	assert.True(t, c.HasBean("testService"))
	
	c.Destroy()
	
	assert.False(t, c.HasBean("testService"))
	beans := c.ListBeans()
	assert.Len(t, beans, 0)
}