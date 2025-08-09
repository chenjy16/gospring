package tests

import (
	"reflect"
	"testing"
	"gospring/annotations"
	"github.com/stretchr/testify/assert"
)

// 测试用的组件
type TestAnnotationService struct {
	name string
	_ string `component:"annotationService" singleton:"true"`
}

func (s *TestAnnotationService) ComponentName() string {
	return "annotationService"
}

func (s *TestAnnotationService) Init() error {
	s.name = "initialized"
	return nil
}

func (s *TestAnnotationService) PostConstruct() error {
	s.name += "_post_construct"
	return nil
}

func (s *TestAnnotationService) Destroy() error {
	s.name = "destroyed"
	return nil
}

func (s *TestAnnotationService) SetBeanName(name string) {
	s.name = name
}

type TestPrototypeService struct {
	_ string `component:"prototypeService" singleton:"false"`
}

func (s *TestPrototypeService) ComponentName() string {
	return "prototypeService"
}

type TestScopeService struct {
	_ string `component:"scopeService" scope:"prototype"`
}

func (s *TestScopeService) ComponentName() string {
	return "scopeService"
}

type TestInitMethodService struct {
	initialized bool
	_ string `component:"initMethodService" init-method:"CustomInit"`
}

func (s *TestInitMethodService) ComponentName() string {
	return "initMethodService"
}

func (s *TestInitMethodService) CustomInit() error {
	s.initialized = true
	return nil
}

type TestDestroyMethodService struct {
	destroyed bool
	_ string `component:"destroyMethodService" destroy-method:"CustomDestroy"`
}

func (s *TestDestroyMethodService) ComponentName() string {
	return "destroyMethodService"
}

func (s *TestDestroyMethodService) CustomDestroy() error {
	s.destroyed = true
	return nil
}

type TestNonComponentService struct {
	name string
}

func TestAnnotationUtils_IsComponent(t *testing.T) {
	utils := annotations.NewAnnotationUtils()
	
	// 测试有component标签的类型
	componentType := reflect.TypeOf(&TestAnnotationService{})
	assert.True(t, utils.IsComponent(componentType))
	
	// 测试没有component标签的类型
	nonComponentType := reflect.TypeOf(&TestNonComponentService{})
	assert.False(t, utils.IsComponent(nonComponentType))
}

func TestAnnotationUtils_GetComponentName(t *testing.T) {
	utils := annotations.NewAnnotationUtils()
	
	// 测试有component标签的类型
	componentType := reflect.TypeOf(&TestAnnotationService{})
	name := utils.GetComponentName(componentType)
	assert.Equal(t, "annotationService", name)
	
	// 测试没有component标签的类型
	nonComponentType := reflect.TypeOf(&TestNonComponentService{})
	name = utils.GetComponentName(nonComponentType)
	assert.Empty(t, name)
}

func TestAnnotationUtils_IsSingleton(t *testing.T) {
	utils := annotations.NewAnnotationUtils()
	
	// 测试单例组件
	singletonType := reflect.TypeOf(&TestAnnotationService{})
	assert.True(t, utils.IsSingleton(singletonType))
	
	// 测试原型组件（singleton="false"）
	prototypeType := reflect.TypeOf(&TestPrototypeService{})
	assert.False(t, utils.IsSingleton(prototypeType))
	
	// 测试原型组件（scope="prototype"）
	scopeType := reflect.TypeOf(&TestScopeService{})
	assert.False(t, utils.IsSingleton(scopeType))
}

func TestAnnotationUtils_GetScope(t *testing.T) {
	utils := annotations.NewAnnotationUtils()
	
	// 测试有scope标签的类型
	scopeType := reflect.TypeOf(&TestScopeService{})
	scope := utils.GetScope(scopeType)
	assert.Equal(t, "prototype", scope)
	
	// 测试没有scope标签的类型
	normalType := reflect.TypeOf(&TestAnnotationService{})
	scope = utils.GetScope(normalType)
	assert.Equal(t, "singleton", scope)
}

func TestAnnotationUtils_HasTag(t *testing.T) {
	utils := annotations.NewAnnotationUtils()
	
	// 测试有component标签的类型
	componentType := reflect.TypeOf(&TestAnnotationService{})
	assert.True(t, utils.HasTag(componentType, "component"))
	assert.True(t, utils.HasTag(componentType, "singleton"))
	
	// 测试没有标签的类型
	nonComponentType := reflect.TypeOf(&TestNonComponentService{})
	assert.False(t, utils.HasTag(nonComponentType, "component"))
}

func TestAnnotationUtils_GetTagValue(t *testing.T) {
	utils := annotations.NewAnnotationUtils()
	
	// 测试获取标签值
	componentType := reflect.TypeOf(&TestAnnotationService{})
	value := utils.GetTagValue(componentType, "_", "component")
	assert.Equal(t, "annotationService", value)
	
	value = utils.GetTagValue(componentType, "_", "singleton")
	assert.Equal(t, "true", value)
}

func TestAnnotationUtils_GetAllTaggedFields(t *testing.T) {
	utils := annotations.NewAnnotationUtils()
	
	// 测试获取所有带component标签的字段
	componentType := reflect.TypeOf(&TestAnnotationService{})
	fields := utils.GetAllTaggedFields(componentType, "component")
	assert.Len(t, fields, 1)
	assert.Equal(t, "_", fields[0].Name)
}

func TestAnnotationUtils_GetInjectFields(t *testing.T) {
	utils := annotations.NewAnnotationUtils()
	
	// 创建一个有inject标签的测试类型
	type TestInjectService struct {
		Service1 interface{} `inject:"service1"`
		Service2 interface{} `inject:"service2"`
		Normal   string
	}
	
	injectType := reflect.TypeOf(&TestInjectService{})
	fields := utils.GetInjectFields(injectType)
	assert.Len(t, fields, 2)
	assert.Equal(t, "Service1", fields[0].Name)
	assert.Equal(t, "Service2", fields[1].Name)
}

func TestInitializer_Interface(t *testing.T) {
	service := &TestAnnotationService{}
	
	// 测试是否实现了Initializer接口
	initializer, ok := interface{}(service).(annotations.Initializer)
	assert.True(t, ok)
	
	// 测试Init方法
	err := initializer.Init()
	assert.NoError(t, err)
	assert.Equal(t, "initialized", service.name)
}

func TestPostConstruct_Interface(t *testing.T) {
	service := &TestAnnotationService{}
	service.name = "test"
	
	// 测试是否实现了PostConstruct接口
	postConstruct, ok := interface{}(service).(annotations.PostConstruct)
	assert.True(t, ok)
	
	// 测试PostConstruct方法
	err := postConstruct.PostConstruct()
	assert.NoError(t, err)
	assert.Equal(t, "test_post_construct", service.name)
}

func TestDestroyer_Interface(t *testing.T) {
	service := &TestAnnotationService{}
	
	// 测试是否实现了Destroyer接口
	destroyer, ok := interface{}(service).(annotations.Destroyer)
	assert.True(t, ok)
	
	// 测试Destroy方法
	err := destroyer.Destroy()
	assert.NoError(t, err)
	assert.Equal(t, "destroyed", service.name)
}

func TestBeanNameAware_Interface(t *testing.T) {
	service := &TestAnnotationService{}
	
	// 测试是否实现了BeanNameAware接口
	beanNameAware, ok := interface{}(service).(annotations.BeanNameAware)
	assert.True(t, ok)
	
	// 测试SetBeanName方法
	beanNameAware.SetBeanName("testBean")
	assert.Equal(t, "testBean", service.name)
}

func TestComponent_Interface(t *testing.T) {
	service := &TestAnnotationService{}
	
	// 测试是否实现了Component接口
	component, ok := interface{}(service).(annotations.Component)
	assert.True(t, ok)
	
	// 测试ComponentName方法
	name := component.ComponentName()
	assert.Equal(t, "annotationService", name)
}

func TestService_Interface(t *testing.T) {
	// Service接口继承自Component接口
	service := &TestAnnotationService{}
	
	// 测试是否实现了Service接口
	serviceInterface, ok := interface{}(service).(annotations.Service)
	assert.True(t, ok)
	
	// 测试继承的ComponentName方法
	name := serviceInterface.ComponentName()
	assert.Equal(t, "annotationService", name)
}

func TestRepository_Interface(t *testing.T) {
	// Repository接口继承自Component接口
	service := &TestAnnotationService{}
	
	// 测试是否实现了Repository接口
	repository, ok := interface{}(service).(annotations.Repository)
	assert.True(t, ok)
	
	// 测试继承的ComponentName方法
	name := repository.ComponentName()
	assert.Equal(t, "annotationService", name)
}

func TestController_Interface(t *testing.T) {
	// Controller接口继承自Component接口
	service := &TestAnnotationService{}
	
	// 测试是否实现了Controller接口
	controller, ok := interface{}(service).(annotations.Controller)
	assert.True(t, ok)
	
	// 测试继承的ComponentName方法
	name := controller.ComponentName()
	assert.Equal(t, "annotationService", name)
}

// 测试边界情况
func TestAnnotationUtils_EdgeCases(t *testing.T) {
	utils := annotations.NewAnnotationUtils()
	
	// 测试nil类型
	assert.False(t, utils.IsComponent(nil))
	assert.Empty(t, utils.GetComponentName(nil))
	assert.True(t, utils.IsSingleton(nil)) // 默认为单例
	
	// 测试非结构体类型
	stringType := reflect.TypeOf("string")
	assert.False(t, utils.IsComponent(stringType))
	assert.Empty(t, utils.GetComponentName(stringType))
}