package annotations

import (
	"reflect"
)

// Component 组件标记接口
type Component interface {
	ComponentName() string
}

// Service 服务层组件标记接口
type Service interface {
	Component
}

// Repository 数据访问层组件标记接口
type Repository interface {
	Component
}

// Controller 控制层组件标记接口
type Controller interface {
	Component
}

// Initializer 初始化接口
type Initializer interface {
	Init() error
}

// Destroyer 销毁接口
type Destroyer interface {
	Destroy() error
}

// PostConstruct 构造后回调接口
type PostConstruct interface {
	PostConstruct() error
}

// PreDestroy 销毁前回调接口
type PreDestroy interface {
	PreDestroy() error
}

// BeanNameAware Bean名称感知接口
type BeanNameAware interface {
	SetBeanName(name string)
}

// ContainerAware 容器感知接口
type ContainerAware interface {
	SetContainer(container interface{})
}

// AnnotationUtils 注解工具类
type AnnotationUtils struct{}

// NewAnnotationUtils 创建注解工具实例
func NewAnnotationUtils() *AnnotationUtils {
	return &AnnotationUtils{}
}

// HasTag 检查结构体是否有指定标签
func (au *AnnotationUtils) HasTag(typ reflect.Type, tagName string) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if _, ok := field.Tag.Lookup(tagName); ok {
			return true
		}
	}
	return false
}

// GetTagValue 获取标签值
func (au *AnnotationUtils) GetTagValue(typ reflect.Type, fieldName, tagName string) string {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	field, ok := typ.FieldByName(fieldName)
	if !ok {
		return ""
	}

	return field.Tag.Get(tagName)
}

// GetAllTaggedFields 获取所有带有指定标签的字段
func (au *AnnotationUtils) GetAllTaggedFields(typ reflect.Type, tagName string) []reflect.StructField {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var fields []reflect.StructField
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if _, ok := field.Tag.Lookup(tagName); ok {
			fields = append(fields, field)
		}
	}
	return fields
}

// IsComponent 检查类型是否为组件
func (au *AnnotationUtils) IsComponent(typ reflect.Type) bool {
	return au.HasTag(typ, "component") || 
		   au.HasTag(typ, "service") || 
		   au.HasTag(typ, "repository") || 
		   au.HasTag(typ, "controller")
}

// GetComponentName 获取组件名称
func (au *AnnotationUtils) GetComponentName(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// 检查各种组件标签
	tags := []string{"component", "service", "repository", "controller"}
	for _, tag := range tags {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if value := field.Tag.Get(tag); value != "" {
				if value == "true" || value == "" {
					return typ.Name()
				}
				return value
			}
		}
	}

	return ""
}

// GetInjectFields 获取需要注入的字段
func (au *AnnotationUtils) GetInjectFields(typ reflect.Type) []reflect.StructField {
	return au.GetAllTaggedFields(typ, "inject")
}

// IsSingleton 检查是否为单例
func (au *AnnotationUtils) IsSingleton(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if value := field.Tag.Get("singleton"); value != "" {
			return value == "true"
		}
		if value := field.Tag.Get("scope"); value != "" {
			return value == "singleton"
		}
	}

	// 默认为单例
	return true
}

// GetScope 获取Bean的作用域
func (au *AnnotationUtils) GetScope(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if value := field.Tag.Get("scope"); value != "" {
			return value
		}
	}

	return "singleton"
}