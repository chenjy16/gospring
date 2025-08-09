package lifecycle

import (
	"fmt"
	"reflect"
	"gospring/annotations"
)

// LifecycleManager 生命周期管理器
type LifecycleManager struct {
	initOrder    []string
	destroyOrder []string
}

// NewLifecycleManager 创建生命周期管理器
func NewLifecycleManager() *LifecycleManager {
	return &LifecycleManager{
		initOrder:    make([]string, 0),
		destroyOrder: make([]string, 0),
	}
}

// ProcessInitialization 处理Bean初始化
func (lm *LifecycleManager) ProcessInitialization(beanName string, instance interface{}) error {
	// 1. 检查是否实现了BeanNameAware接口
	if aware, ok := instance.(annotations.BeanNameAware); ok {
		aware.SetBeanName(beanName)
	}

	// 2. 检查是否实现了Initializer接口
	if initializer, ok := instance.(annotations.Initializer); ok {
		if err := initializer.Init(); err != nil {
			return fmt.Errorf("failed to initialize bean '%s': %v", beanName, err)
		}
	}

	// 3. 检查是否实现了PostConstruct接口
	if postConstruct, ok := instance.(annotations.PostConstruct); ok {
		if err := postConstruct.PostConstruct(); err != nil {
			return fmt.Errorf("failed to execute post construct for bean '%s': %v", beanName, err)
		}
	}

	// 4. 调用自定义初始化方法（通过反射）
	if err := lm.callInitMethod(instance); err != nil {
		return fmt.Errorf("failed to call init method for bean '%s': %v", beanName, err)
	}

	// 记录初始化顺序
	lm.initOrder = append(lm.initOrder, beanName)

	return nil
}

// ProcessDestruction 处理Bean销毁
func (lm *LifecycleManager) ProcessDestruction(beanName string, instance interface{}) error {
	// 1. 检查是否实现了PreDestroy接口
	if preDestroy, ok := instance.(annotations.PreDestroy); ok {
		if err := preDestroy.PreDestroy(); err != nil {
			return fmt.Errorf("failed to execute pre destroy for bean '%s': %v", beanName, err)
		}
	}

	// 2. 检查是否实现了Destroyer接口
	if destroyer, ok := instance.(annotations.Destroyer); ok {
		if err := destroyer.Destroy(); err != nil {
			return fmt.Errorf("failed to destroy bean '%s': %v", beanName, err)
		}
	}

	// 3. 调用自定义销毁方法（通过反射）
	if err := lm.callDestroyMethod(instance); err != nil {
		return fmt.Errorf("failed to call destroy method for bean '%s': %v", beanName, err)
	}

	// 记录销毁顺序（逆序）
	lm.destroyOrder = append([]string{beanName}, lm.destroyOrder...)

	return nil
}

// callInitMethod 通过反射调用初始化方法
func (lm *LifecycleManager) callInitMethod(instance interface{}) error {
	val := reflect.ValueOf(instance)
	typ := reflect.TypeOf(instance)

	// 查找init方法
	initMethods := []string{"Init", "Initialize", "PostConstruct", "AfterPropertiesSet"}
	
	for _, methodName := range initMethods {
		method := val.MethodByName(methodName)
		if method.IsValid() && method.Type().NumIn() == 0 {
			// 调用无参数的初始化方法
			results := method.Call(nil)
			
			// 检查返回值是否有错误
			if len(results) > 0 {
				if err, ok := results[0].Interface().(error); ok && err != nil {
					return err
				}
			}
			break
		}
	}

	// 检查结构体标签中的初始化方法
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if initMethod := field.Tag.Get("init-method"); initMethod != "" {
			method := val.MethodByName(initMethod)
			if method.IsValid() {
				results := method.Call(nil)
				if len(results) > 0 {
					if err, ok := results[0].Interface().(error); ok && err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// callDestroyMethod 通过反射调用销毁方法
func (lm *LifecycleManager) callDestroyMethod(instance interface{}) error {
	val := reflect.ValueOf(instance)
	typ := reflect.TypeOf(instance)

	// 查找destroy方法
	destroyMethods := []string{"Destroy", "Close", "Cleanup", "PreDestroy"}
	
	for _, methodName := range destroyMethods {
		method := val.MethodByName(methodName)
		if method.IsValid() && method.Type().NumIn() == 0 {
			// 调用无参数的销毁方法
			results := method.Call(nil)
			
			// 检查返回值是否有错误
			if len(results) > 0 {
				if err, ok := results[0].Interface().(error); ok && err != nil {
					return err
				}
			}
			break
		}
	}

	// 检查结构体标签中的销毁方法
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if destroyMethod := field.Tag.Get("destroy-method"); destroyMethod != "" {
			method := val.MethodByName(destroyMethod)
			if method.IsValid() {
				results := method.Call(nil)
				if len(results) > 0 {
					if err, ok := results[0].Interface().(error); ok && err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// GetInitOrder 获取初始化顺序
func (lm *LifecycleManager) GetInitOrder() []string {
	return lm.initOrder
}

// GetDestroyOrder 获取销毁顺序
func (lm *LifecycleManager) GetDestroyOrder() []string {
	return lm.destroyOrder
}

// Reset 重置生命周期管理器
func (lm *LifecycleManager) Reset() {
	lm.initOrder = make([]string, 0)
	lm.destroyOrder = make([]string, 0)
}