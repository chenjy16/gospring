package context

import (
	"fmt"
	"reflect"
	"gospring/container"
	"gospring/scanner"
	"gospring/lifecycle"
	"gospring/annotations"
)

// ApplicationContext 应用上下文
type ApplicationContext struct {
	container         *container.Container
	scanner           *scanner.ComponentScanner
	lifecycleManager  *lifecycle.LifecycleManager
	annotationUtils   *annotations.AnnotationUtils
	started           bool
}

// NewApplicationContext 创建新的应用上下文
func NewApplicationContext() *ApplicationContext {
	c := container.NewContainer()
	return &ApplicationContext{
		container:        c,
		scanner:          scanner.NewComponentScanner(c),
		lifecycleManager: lifecycle.NewLifecycleManager(),
		annotationUtils:  annotations.NewAnnotationUtils(),
		started:          false,
	}
}

// RegisterBean 注册Bean
func (ctx *ApplicationContext) RegisterBean(name string, instance interface{}) error {
	// 检查是否为单例
	typ := reflect.TypeOf(instance)
	singleton := ctx.annotationUtils.IsSingleton(typ)

	var err error
	if singleton {
		err = ctx.container.RegisterSingleton(name, instance)
	} else {
		err = ctx.container.RegisterPrototype(name, instance)
	}

	if err != nil {
		return err
	}

	// 如果上下文已启动，立即处理生命周期
	if ctx.started {
		return ctx.lifecycleManager.ProcessInitialization(name, instance)
	}

	return nil
}

// RegisterComponent 注册组件
func (ctx *ApplicationContext) RegisterComponent(instance interface{}) error {
	return ctx.scanner.ScanComponent(instance)
}

// RegisterComponents 批量注册组件
func (ctx *ApplicationContext) RegisterComponents(components ...interface{}) error {
	return ctx.scanner.ScanAndRegister(components...)
}

// RegisterByInterface 根据接口注册实现
func (ctx *ApplicationContext) RegisterByInterface(interfaceType reflect.Type, implementation interface{}, name string) error {
	return ctx.scanner.RegisterWithInterface(interfaceType, implementation, name)
}

// GetBean 获取Bean
func (ctx *ApplicationContext) GetBean(name string) interface{} {
	return ctx.container.GetBean(name)
}

// GetBeanByType 根据类型获取Bean
func (ctx *ApplicationContext) GetBeanByType(typ reflect.Type) interface{} {
	return ctx.container.GetBeanByType(typ)
}

// GetBeanT 泛型方式获取Bean（Go 1.18+）
func GetBeanT[T any](ctx *ApplicationContext, name string) T {
	var zero T
	bean := ctx.GetBean(name)
	if bean == nil {
		return zero
	}
	
	if result, ok := bean.(T); ok {
		return result
	}
	
	return zero
}

// Start 启动应用上下文
func (ctx *ApplicationContext) Start() error {
	if ctx.started {
		return fmt.Errorf("application context is already started")
	}

	// 1. 执行依赖注入
	if err := ctx.container.WireAll(); err != nil {
		return fmt.Errorf("failed to wire dependencies: %v", err)
	}

	// 2. 处理所有Bean的生命周期初始化
	beanNames := ctx.container.ListBeans()
	for _, beanName := range beanNames {
		bean := ctx.container.GetBean(beanName)
		if bean != nil {
			if err := ctx.lifecycleManager.ProcessInitialization(beanName, bean); err != nil {
				return fmt.Errorf("failed to initialize bean '%s': %v", beanName, err)
			}
		}
	}

	ctx.started = true
	return nil
}

// Stop 停止应用上下文
func (ctx *ApplicationContext) Stop() error {
	if !ctx.started {
		return fmt.Errorf("application context is not started")
	}

	// 按逆序销毁Bean
	beanNames := ctx.container.ListBeans()
	for i := len(beanNames) - 1; i >= 0; i-- {
		beanName := beanNames[i]
		bean := ctx.container.GetBean(beanName)
		if bean != nil {
			if err := ctx.lifecycleManager.ProcessDestruction(beanName, bean); err != nil {
				// 记录错误但继续销毁其他Bean
				fmt.Printf("Error destroying bean '%s': %v\n", beanName, err)
			}
		}
	}

	// 销毁容器
	ctx.container.Destroy()
	ctx.started = false

	return nil
}

// Refresh 刷新上下文
func (ctx *ApplicationContext) Refresh() error {
	if ctx.started {
		if err := ctx.Stop(); err != nil {
			return err
		}
	}
	return ctx.Start()
}

// IsStarted 检查上下文是否已启动
func (ctx *ApplicationContext) IsStarted() bool {
	return ctx.started
}

// HasBean 检查是否存在指定Bean
func (ctx *ApplicationContext) HasBean(name string) bool {
	return ctx.container.HasBean(name)
}

// ListBeans 列出所有Bean名称
func (ctx *ApplicationContext) ListBeans() []string {
	return ctx.container.ListBeans()
}

// GetBeanDefinition 获取Bean定义
func (ctx *ApplicationContext) GetBeanDefinition(name string) *container.BeanDefinition {
	return ctx.container.GetBeanDefinition(name)
}

// GetContainer 获取底层容器
func (ctx *ApplicationContext) GetContainer() *container.Container {
	return ctx.container
}

// GetLifecycleManager 获取生命周期管理器
func (ctx *ApplicationContext) GetLifecycleManager() *lifecycle.LifecycleManager {
	return ctx.lifecycleManager
}

// AutoWire 自动装配指定实例的依赖
func (ctx *ApplicationContext) AutoWire(instance interface{}) error {
	return ctx.container.InjectDependencies(instance)
}

// CreateBean 创建并注册新Bean
func (ctx *ApplicationContext) CreateBean(name string, factory func() interface{}) error {
	instance := factory()
	return ctx.RegisterBean(name, instance)
}

// GetBeansOfType 获取指定类型的所有Bean
func (ctx *ApplicationContext) GetBeansOfType(typ reflect.Type) map[string]interface{} {
	result := make(map[string]interface{})
	beanNames := ctx.container.ListBeans()
	
	for _, beanName := range beanNames {
		bean := ctx.container.GetBean(beanName)
		if bean != nil {
			beanType := reflect.TypeOf(bean)
			if beanType.AssignableTo(typ) || beanType.Implements(typ) {
				result[beanName] = bean
			}
		}
	}
	
	return result
}