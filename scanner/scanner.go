package scanner

import (
	"fmt"
	"reflect"
	"strings"
	"gospring/container"
)

// ComponentScanner 组件扫描器
type ComponentScanner struct {
	container *container.Container
	packages  []string
}

// NewComponentScanner 创建新的组件扫描器
func NewComponentScanner(c *container.Container) *ComponentScanner {
	return &ComponentScanner{
		container: c,
		packages:  make([]string, 0),
	}
}

// AddPackage 添加要扫描的包
func (s *ComponentScanner) AddPackage(pkg string) {
	s.packages = append(s.packages, pkg)
}

// ScanComponent 扫描并注册组件
func (s *ComponentScanner) ScanComponent(instance interface{}) error {
	typ := reflect.TypeOf(instance)
	val := reflect.ValueOf(instance)

	// 如果是指针，获取其指向的类型
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	// 检查是否有component标签
	componentName := s.getComponentName(typ)
	if componentName == "" {
		return fmt.Errorf("type %v is not a component", typ)
	}

	// 检查是否为单例
	singleton := s.isSingleton(typ)

	// 注册到容器
	if singleton {
		return s.container.RegisterSingleton(componentName, instance)
	} else {
		return s.container.RegisterPrototype(componentName, instance)
	}
}

// getComponentName 获取组件名称
func (s *ComponentScanner) getComponentName(typ reflect.Type) string {
	// 检查结构体标签
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if componentTag := field.Tag.Get("component"); componentTag != "" {
			if componentTag == "true" || componentTag == "" {
				// 使用类型名称作为组件名
				return strings.ToLower(typ.Name())
			}
			return componentTag
		}
	}

	// 检查类型是否有特定的命名约定
	typeName := typ.Name()
	if strings.HasSuffix(typeName, "Service") || 
	   strings.HasSuffix(typeName, "Repository") || 
	   strings.HasSuffix(typeName, "Controller") ||
	   strings.HasSuffix(typeName, "Component") {
		return strings.ToLower(typeName)
	}

	return ""
}

// isSingleton 检查是否为单例
func (s *ComponentScanner) isSingleton(typ reflect.Type) bool {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if singletonTag := field.Tag.Get("singleton"); singletonTag != "" {
			return singletonTag == "true"
		}
	}
	
	// 默认为单例
	return true
}

// ScanAndRegister 扫描多个组件并注册
func (s *ComponentScanner) ScanAndRegister(components ...interface{}) error {
	for _, component := range components {
		if err := s.ScanComponent(component); err != nil {
			return fmt.Errorf("failed to scan component %T: %v", component, err)
		}
	}
	return nil
}

// AutoScan 自动扫描结构体字段中的组件标签
func (s *ComponentScanner) AutoScan(instance interface{}) error {
	typ := reflect.TypeOf(instance)
	val := reflect.ValueOf(instance)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	// 扫描所有字段
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// 检查是否有component标签
		if componentTag := field.Tag.Get("component"); componentTag != "" {
			if fieldVal.IsValid() && !fieldVal.IsNil() {
				if err := s.ScanComponent(fieldVal.Interface()); err != nil {
					return err
				}
			}
		}

		// 递归扫描嵌套结构体
		if field.Type.Kind() == reflect.Struct {
			if fieldVal.CanInterface() {
				if err := s.AutoScan(fieldVal.Addr().Interface()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// RegisterWithInterface 注册实现了特定接口的组件
func (s *ComponentScanner) RegisterWithInterface(interfaceType reflect.Type, implementation interface{}, name string) error {
	return s.container.RegisterByInterface(interfaceType, implementation, name)
}

// ScanPackageComponents 扫描包中的组件（模拟实现）
func (s *ComponentScanner) ScanPackageComponents() error {
	// 这里可以实现更复杂的包扫描逻辑
	// 由于Go的反射限制，实际项目中可能需要使用代码生成或其他方式
	fmt.Println("Package scanning is not fully implemented in this demo")
	return nil
}