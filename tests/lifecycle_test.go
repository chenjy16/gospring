package tests

import (
	"errors"
	"testing"
	"gospring/lifecycle"
	"gospring/logging"
	"github.com/stretchr/testify/assert"
)

// 测试用的组件
type TestLifecycleService struct {
	name         string
	initialized  bool
	postConstruct bool
	preDestroyed bool
	destroyed    bool
	beanName     string
	initError    error
	destroyError error
}

func (s *TestLifecycleService) Init() error {
	if s.initError != nil {
		return s.initError
	}
	s.initialized = true
	s.name = "initialized"
	return nil
}

func (s *TestLifecycleService) PostConstruct() error {
	s.postConstruct = true
	s.name += "_post_construct"
	return nil
}

func (s *TestLifecycleService) PreDestroy() error {
	s.preDestroyed = true
	return nil
}

func (s *TestLifecycleService) Destroy() error {
	if s.destroyError != nil {
		return s.destroyError
	}
	s.destroyed = true
	s.name = "destroyed"
	return nil
}

func (s *TestLifecycleService) SetBeanName(name string) {
	s.beanName = name
}

type TestCustomMethodService struct {
	customInitialized bool
	customDestroyed   bool
	_ string `init-method:"CustomInit" destroy-method:"CustomDestroy"`
}

func (s *TestCustomMethodService) CustomInit() error {
	s.customInitialized = true
	return nil
}

func (s *TestCustomMethodService) CustomDestroy() error {
	s.customDestroyed = true
	return nil
}

type TestErrorService struct {
	shouldFailInit    bool
	shouldFailDestroy bool
}

func (s *TestErrorService) Init() error {
	if s.shouldFailInit {
		return errors.New("init failed")
	}
	return nil
}

func (s *TestErrorService) Destroy() error {
	if s.shouldFailDestroy {
		return errors.New("destroy failed")
	}
	return nil
}

type TestReflectionMethodService struct {
	initialized bool
	destroyed   bool
}

func (s *TestReflectionMethodService) Initialize() error {
	s.initialized = true
	return nil
}

func (s *TestReflectionMethodService) Close() error {
	s.destroyed = true
	return nil
}

func TestNewLifecycleManager(t *testing.T) {
	lm := lifecycle.NewLifecycleManager()
	assert.NotNil(t, lm)
	assert.NotNil(t, lm.GetLogger())
	assert.Empty(t, lm.GetInitOrder())
	assert.Empty(t, lm.GetDestroyOrder())
}

func TestNewLifecycleManagerWithLogger(t *testing.T) {
	logger := logging.NopLogger
	lm := lifecycle.NewLifecycleManagerWithLogger(logger)
	assert.NotNil(t, lm)
	assert.Equal(t, logger, lm.GetLogger())
}

func TestLifecycleManager_SetGetLogger(t *testing.T) {
	lm := lifecycle.NewLifecycleManager()
	logger := logging.NopLogger
	
	lm.SetLogger(logger)
	assert.Equal(t, logger, lm.GetLogger())
}

func TestLifecycleManager_ProcessInitialization_Success(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestLifecycleService{}
	
	err := lm.ProcessInitialization("testService", service)
	
	assert.NoError(t, err)
	assert.True(t, service.initialized)
	assert.True(t, service.postConstruct)
	assert.Equal(t, "testService", service.beanName)
	assert.Equal(t, "initialized_post_construct", service.name)
	assert.Contains(t, lm.GetInitOrder(), "testService")
}

func TestLifecycleManager_ProcessInitialization_InitError(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestLifecycleService{
		initError: errors.New("init failed"),
	}
	
	err := lm.ProcessInitialization("testService", service)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to initialize bean 'testService'")
	assert.False(t, service.postConstruct) // PostConstruct 不应该被调用
	assert.NotContains(t, lm.GetInitOrder(), "testService")
}

func TestLifecycleManager_ProcessInitialization_CustomMethod(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestCustomMethodService{}
	
	err := lm.ProcessInitialization("customService", service)
	
	assert.NoError(t, err)
	assert.True(t, service.customInitialized)
	assert.Contains(t, lm.GetInitOrder(), "customService")
}

func TestLifecycleManager_ProcessInitialization_ReflectionMethod(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestReflectionMethodService{}
	
	err := lm.ProcessInitialization("reflectionService", service)
	
	assert.NoError(t, err)
	assert.True(t, service.initialized)
	assert.Contains(t, lm.GetInitOrder(), "reflectionService")
}

func TestLifecycleManager_ProcessDestruction_Success(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestLifecycleService{}
	
	err := lm.ProcessDestruction("testService", service)
	
	assert.NoError(t, err)
	assert.True(t, service.preDestroyed)
	assert.True(t, service.destroyed)
	assert.Equal(t, "destroyed", service.name)
	assert.Contains(t, lm.GetDestroyOrder(), "testService")
}

func TestLifecycleManager_ProcessDestruction_Error(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestLifecycleService{
		destroyError: errors.New("destroy failed"),
	}
	
	err := lm.ProcessDestruction("testService", service)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to destroy bean 'testService'")
	assert.True(t, service.preDestroyed) // PreDestroy 应该被调用
	assert.Contains(t, lm.GetDestroyOrder(), "testService")
}

func TestLifecycleManager_ProcessDestruction_CustomMethod(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestCustomMethodService{}
	
	err := lm.ProcessDestruction("customService", service)
	
	assert.NoError(t, err)
	assert.True(t, service.customDestroyed)
	assert.Contains(t, lm.GetDestroyOrder(), "customService")
}

func TestLifecycleManager_ProcessDestruction_ReflectionMethod(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestReflectionMethodService{}
	
	err := lm.ProcessDestruction("reflectionService", service)
	
	assert.NoError(t, err)
	assert.True(t, service.destroyed)
	assert.Contains(t, lm.GetDestroyOrder(), "reflectionService")
}

func TestLifecycleManager_InitDestroyOrder(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	
	service1 := &TestLifecycleService{}
	service2 := &TestLifecycleService{}
	service3 := &TestLifecycleService{}
	
	// 初始化顺序
	lm.ProcessInitialization("service1", service1)
	lm.ProcessInitialization("service2", service2)
	lm.ProcessInitialization("service3", service3)
	
	initOrder := lm.GetInitOrder()
	assert.Equal(t, []string{"service1", "service2", "service3"}, initOrder)
	
	// 销毁顺序（逆序）
	lm.ProcessDestruction("service1", service1)
	lm.ProcessDestruction("service2", service2)
	lm.ProcessDestruction("service3", service3)
	
	destroyOrder := lm.GetDestroyOrder()
	assert.Equal(t, []string{"service3", "service2", "service1"}, destroyOrder)
}

func TestLifecycleManager_Reset(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestLifecycleService{}
	
	lm.ProcessInitialization("testService", service)
	lm.ProcessDestruction("testService", service)
	
	assert.NotEmpty(t, lm.GetInitOrder())
	assert.NotEmpty(t, lm.GetDestroyOrder())
	
	lm.Reset()
	
	assert.Empty(t, lm.GetInitOrder())
	assert.Empty(t, lm.GetDestroyOrder())
}

func TestLifecycleManager_BeanNameAware(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	service := &TestLifecycleService{}
	
	err := lm.ProcessInitialization("myBean", service)
	
	assert.NoError(t, err)
	assert.Equal(t, "myBean", service.beanName)
}

func TestLifecycleManager_NonLifecycleBean(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	
	// 测试没有实现任何生命周期接口的Bean
	type SimpleBean struct {
		value string
	}
	
	bean := &SimpleBean{value: "test"}
	
	err := lm.ProcessInitialization("simpleBean", bean)
	assert.NoError(t, err)
	assert.Contains(t, lm.GetInitOrder(), "simpleBean")
	
	err = lm.ProcessDestruction("simpleBean", bean)
	assert.NoError(t, err)
	assert.Contains(t, lm.GetDestroyOrder(), "simpleBean")
}

func TestLifecycleManager_WithLogging(t *testing.T) {
	// 使用控制台日志器测试日志记录
	logger := logging.NewConsoleLogger()
	lm := lifecycle.NewLifecycleManagerWithLogger(logger)
	service := &TestLifecycleService{}
	
	// 测试初始化日志
	err := lm.ProcessInitialization("loggedService", service)
	assert.NoError(t, err)
	
	// 测试销毁日志
	err = lm.ProcessDestruction("loggedService", service)
	assert.NoError(t, err)
}

func TestLifecycleManager_ErrorHandling(t *testing.T) {
	lm := lifecycle.NewLifecycleManagerWithLogger(logging.NopLogger)
	
	// 测试初始化错误
	errorService := &TestErrorService{shouldFailInit: true}
	err := lm.ProcessInitialization("errorService", errorService)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "init failed")
	
	// 测试销毁错误
	errorService = &TestErrorService{shouldFailDestroy: true}
	err = lm.ProcessDestruction("errorService", errorService)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "destroy failed")
}