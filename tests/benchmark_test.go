package tests

import (
	"reflect"
	"testing"

	"gospring/container"
	"gospring/context"
)

// BenchmarkService 用于性能测试的服务
type BenchmarkService struct {
	_ string `component:"benchmarkService"`
}

func (s *BenchmarkService) DoWork() string {
	return "work done"
}

// BenchmarkRepository 用于性能测试的仓库
type BenchmarkRepository struct {
	_ string `component:"benchmarkRepository"`
}

func (r *BenchmarkRepository) GetData() string {
	return "data"
}

// BenchmarkController 用于性能测试的控制器
type BenchmarkController struct {
	Service    *BenchmarkService    `inject:""`
	Repository *BenchmarkRepository `inject:""`
	_          string               `component:"benchmarkController"`
}

func (c *BenchmarkController) Handle() string {
	return c.Service.DoWork() + " " + c.Repository.GetData()
}

// BenchmarkContainerRegister 测试容器注册性能
func BenchmarkContainerRegister(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := container.NewContainer()
		c.RegisterSingleton("service", &BenchmarkService{})
		c.RegisterSingleton("repository", &BenchmarkRepository{})
		c.RegisterSingleton("controller", &BenchmarkController{})
	}
}

// BenchmarkContainerGetBean 测试获取Bean性能
func BenchmarkContainerGetBean(b *testing.B) {
	c := container.NewContainer()
	c.RegisterSingleton("service", &BenchmarkService{})
	c.RegisterSingleton("repository", &BenchmarkRepository{})
	c.RegisterSingleton("controller", &BenchmarkController{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetBean("service")
		c.GetBean("repository")
		c.GetBean("controller")
	}
}

// BenchmarkContainerGetBeanByType 测试按类型获取Bean性能
func BenchmarkContainerGetBeanByType(b *testing.B) {
	c := container.NewContainer()
	c.RegisterSingleton("service", &BenchmarkService{})
	c.RegisterSingleton("repository", &BenchmarkRepository{})
	c.RegisterSingleton("controller", &BenchmarkController{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetBeanByType(reflect.TypeOf(&BenchmarkService{}))
		c.GetBeanByType(reflect.TypeOf(&BenchmarkRepository{}))
		c.GetBeanByType(reflect.TypeOf(&BenchmarkController{}))
	}
}

// BenchmarkDependencyInjection 测试依赖注入性能
func BenchmarkDependencyInjection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.NewApplicationContext()
		ctx.RegisterComponent(&BenchmarkService{})
		ctx.RegisterComponent(&BenchmarkRepository{})
		ctx.RegisterComponent(&BenchmarkController{})
		ctx.Start()
		ctx.Stop()
	}
}

// BenchmarkApplicationContextStart 测试应用上下文启动性能
func BenchmarkApplicationContextStart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.NewApplicationContext()
		ctx.RegisterComponent(&BenchmarkService{})
		ctx.RegisterComponent(&BenchmarkRepository{})
		ctx.RegisterComponent(&BenchmarkController{})
		
		b.StartTimer()
		ctx.Start()
		b.StopTimer()
		
		ctx.Stop()
	}
}

// BenchmarkPrototypeCreation 测试原型Bean创建性能
func BenchmarkPrototypeCreation(b *testing.B) {
	c := container.NewContainer()
	c.RegisterPrototype("service", &BenchmarkService{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetBean("service")
	}
}

// BenchmarkConcurrentAccess 测试并发访问性能
func BenchmarkConcurrentAccess(b *testing.B) {
	c := container.NewContainer()
	c.RegisterSingleton("service", &BenchmarkService{})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.GetBean("service")
		}
	})
}