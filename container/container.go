package container

import (
	"fmt"
	"reflect"
	"sync"
	"time"
	"gospring/logging"
)

// BeanDefinition 定义Bean的元数据
type BeanDefinition struct {
	Name      string
	Type      reflect.Type
	Value     reflect.Value
	Singleton bool
	Instance  interface{}
	mutex     sync.RWMutex
}

// Container IoC容器
type Container struct {
	beans       map[string]*BeanDefinition
	typeMapping map[reflect.Type]string // 类型到Bean名称的映射
	mutex       sync.RWMutex
	logger      logging.Logger // 日志器
}

// NewContainer 创建新的容器实例
func NewContainer() *Container {
	return NewContainerWithLogger(logging.NewConsoleLogger())
}

// NewContainerWithLogger 创建带有指定日志器的容器实例
func NewContainerWithLogger(logger logging.Logger) *Container {
	container := &Container{
		beans:       make(map[string]*BeanDefinition),
		typeMapping: make(map[reflect.Type]string),
		logger:      logger,
	}
	
	// 记录容器创建事件
	container.logger.LogEvent(&logging.ContainerCreated{
		Timestamp: time.Now(),
	})
	
	return container
}

// RegisterSingleton 注册单例Bean
func (c *Container) RegisterSingleton(name string, instance interface{}) error {
	return c.registerBean(name, instance, true)
}

// RegisterPrototype 注册原型Bean
func (c *Container) RegisterPrototype(name string, instance interface{}) error {
	return c.registerBean(name, instance, false)
}

// registerBean 内部注册Bean方法
func (c *Container) registerBean(name string, instance interface{}, singleton bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.beans[name]; exists {
		return fmt.Errorf("bean with name '%s' already exists", name)
	}

	val := reflect.ValueOf(instance)
	typ := reflect.TypeOf(instance)
	originalType := typ

	// 如果是指针，获取其指向的类型
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	beanDef := &BeanDefinition{
		Name:      name,
		Type:      typ,
		Value:     val,
		Singleton: singleton,
		Instance:  instance,
	}

	c.beans[name] = beanDef
	// 同时注册指针类型和元素类型的映射
	c.typeMapping[typ] = name
	c.typeMapping[originalType] = name

	// 如果实现了接口，也注册接口映射
	c.registerInterfaces(instance, name)

	// 记录组件注册事件
	scope := "singleton"
	if !singleton {
		scope = "prototype"
	}
	c.logger.LogEvent(&logging.ComponentRegistered{
		Timestamp:     time.Now(),
		ComponentID:   name,
		ComponentType: typ.String(),
		Scope:         scope,
	})

	return nil
}

// registerInterfaces 注册接口映射
func (c *Container) registerInterfaces(instance interface{}, beanName string) {
	typ := reflect.TypeOf(instance)
	
	// 遍历所有接口
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		// 检查是否实现了某个接口
		if method.Type.NumIn() > 0 {
			// 这里可以扩展更复杂的接口检测逻辑
		}
	}
}

// GetBean 获取Bean实例
func (c *Container) GetBean(name string) interface{} {
	c.mutex.RLock()
	beanDef, exists := c.beans[name]
	c.mutex.RUnlock()

	if !exists {
		return nil
	}

	if beanDef.Singleton {
		return beanDef.Instance
	}

	// 原型模式，创建新实例
	return c.createNewInstance(beanDef)
}

// GetBeanByType 根据类型获取Bean
func (c *Container) GetBeanByType(typ reflect.Type) interface{} {
	c.mutex.RLock()
	beanName, exists := c.typeMapping[typ]
	c.mutex.RUnlock()

	if !exists {
		return nil
	}

	return c.GetBean(beanName)
}

// createNewInstance 创建新的实例（用于原型模式）
func (c *Container) createNewInstance(beanDef *BeanDefinition) interface{} {
	start := time.Now()
	
	// 创建新实例
	newVal := reflect.New(beanDef.Type)
	newInstance := newVal.Interface()

	// 执行依赖注入
	c.InjectDependencies(newInstance)

	// 记录组件创建事件
	c.logger.LogEvent(&logging.ComponentCreated{
		Timestamp:     time.Now(),
		ComponentID:   beanDef.Name,
		ComponentType: beanDef.Type.String(),
		CreationTime:  time.Since(start),
	})

	return newInstance
}

// InjectDependencies 执行依赖注入
func (c *Container) InjectDependencies(instance interface{}) error {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	// 遍历所有字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 检查inject标签
		injectTag := fieldType.Tag.Get("inject")
		if injectTag == "" {
			continue
		}

		// 如果字段不可设置，跳过
		if !field.CanSet() {
			continue
		}

		var dependency interface{}

		// 如果标签指定了Bean名称
		if injectTag != "" && injectTag != "true" {
			dependency = c.GetBean(injectTag)
		} else {
			// 根据类型查找
			dependency = c.GetBeanByType(fieldType.Type)
		}

		if dependency != nil {
			depVal := reflect.ValueOf(dependency)
			if depVal.Type().AssignableTo(field.Type()) {
				field.Set(depVal)
				
				// 记录依赖注入成功事件
				c.logger.LogEvent(&logging.DependencyInjected{
					Timestamp:      time.Now(),
					TargetType:     typ.String(),
					DependencyType: depVal.Type().String(),
					FieldName:      fieldType.Name,
					ByType:         injectTag == "" || injectTag == "true",
					ByName:         injectTag != "" && injectTag != "true",
				})
			}
		} else {
			// 记录依赖注入失败事件
			c.logger.LogEvent(&logging.DependencyInjectionFailed{
				Timestamp:      time.Now(),
				TargetType:     typ.String(),
				DependencyType: fieldType.Type.String(),
				FieldName:      fieldType.Name,
				Error:          fmt.Errorf("dependency not found"),
			})
		}
	}

	return nil
}

// WireAll 对所有已注册的Bean执行依赖注入
func (c *Container) WireAll() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, beanDef := range c.beans {
		if err := c.InjectDependencies(beanDef.Instance); err != nil {
			return fmt.Errorf("failed to inject dependencies for bean '%s': %v", beanDef.Name, err)
		}
	}

	return nil
}

// ListBeans 列出所有注册的Bean
func (c *Container) ListBeans() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var names []string
	for name := range c.beans {
		names = append(names, name)
	}

	return names
}

// HasBean 检查是否存在指定名称的Bean
func (c *Container) HasBean(name string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	_, exists := c.beans[name]
	return exists
}

// GetBeanDefinition 获取Bean定义
func (c *Container) GetBeanDefinition(name string) *BeanDefinition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.beans[name]
}

// RegisterByInterface 根据接口注册实现
func (c *Container) RegisterByInterface(interfaceType reflect.Type, implementation interface{}, name string) error {
	implType := reflect.TypeOf(implementation)
	
	// 检查是否实现了接口
	if !implType.Implements(interfaceType) {
		return fmt.Errorf("type %v does not implement interface %v", implType, interfaceType)
	}

	// 注册实现
	if err := c.RegisterSingleton(name, implementation); err != nil {
		return err
	}

	// 注册接口映射
	c.mutex.Lock()
	c.typeMapping[interfaceType] = name
	c.mutex.Unlock()

	return nil
}

// SetLogger 设置容器的日志器
func (c *Container) SetLogger(logger logging.Logger) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.logger = logger
}

// GetLogger 获取容器的日志器
func (c *Container) GetLogger() logging.Logger {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.logger
}

// Destroy 销毁容器，清理资源
func (c *Container) Destroy() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 调用所有Bean的销毁方法（如果有的话）
	for _, beanDef := range c.beans {
		if destroyer, ok := beanDef.Instance.(interface{ Destroy() }); ok {
			destroyer.Destroy()
		}
		
		// 记录组件销毁事件
		c.logger.LogEvent(&logging.ComponentDestroyed{
			Timestamp:     time.Now(),
			ComponentID:   beanDef.Name,
			ComponentType: beanDef.Type.String(),
		})
	}

	// 清理映射
	c.beans = make(map[string]*BeanDefinition)
	c.typeMapping = make(map[reflect.Type]string)
}