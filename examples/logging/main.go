package main

import (
	"fmt"
	"log"
	"os"
	"gospring/context"
	"gospring/logging"
)

// UserService 用户服务
type UserService struct {
	Name string `component:"userService"`
}

func (u *UserService) Init() error {
	fmt.Println("UserService initializing...")
	return nil
}

func (u *UserService) PostConstruct() error {
	fmt.Println("UserService post construct...")
	return nil
}

func (u *UserService) GetUser(id int) string {
	return fmt.Sprintf("User-%d", id)
}

func (u *UserService) PreDestroy() error {
	fmt.Println("UserService pre destroy...")
	return nil
}

func (u *UserService) Destroy() error {
	fmt.Println("UserService destroying...")
	return nil
}

// OrderService 订单服务
type OrderService struct {
	UserService *UserService `autowired:"true"`
	Name        string       `component:"orderService"`
}

func (o *OrderService) Init() error {
	fmt.Println("OrderService initializing...")
	return nil
}

func (o *OrderService) CreateOrder(userId int) string {
	user := o.UserService.GetUser(userId)
	return fmt.Sprintf("Order for %s", user)
}

func (o *OrderService) Destroy() error {
	fmt.Println("OrderService destroying...")
	return nil
}

func main() {
	fmt.Println("=== GoSpring Logging Demo ===")

	// 创建不同类型的日志器
	consoleLogger := logging.NewConsoleLogger()
	standardLogger := logging.NewStandardLogger(log.New(os.Stdout, "[GoSpring] ", log.LstdFlags))
	
	// 创建过滤日志器，只记录错误和依赖注入事件
	filteredLogger := logging.NewFilteredLogger(consoleLogger, func(event logging.Event) bool {
		switch event.(type) {
		case *logging.DependencyInjected, *logging.DependencyInjectionFailed:
			return true
		case *logging.ComponentCreated, *logging.ComponentDestroyed:
			return true
		default:
			return false
		}
	})

	// 创建多重日志器
	multiLogger := logging.NewMultiLogger(consoleLogger, standardLogger)

	fmt.Println("\n--- 使用控制台日志器 ---")
	demoWithLogger(consoleLogger)

	fmt.Println("\n--- 使用过滤日志器（只显示组件和依赖注入事件）---")
	demoWithLogger(filteredLogger)

	fmt.Println("\n--- 使用多重日志器 ---")
	demoWithLogger(multiLogger)
}

func demoWithLogger(logger logging.Logger) {
	// 创建带有指定日志器的应用上下文
	ctx := context.NewApplicationContextWithLogger(logger)

	// 创建服务实例
	userService := &UserService{Name: "UserService"}
	orderService := &OrderService{Name: "OrderService"}

	// 注册组件
	ctx.RegisterComponent(userService)
	ctx.RegisterComponent(orderService)

	// 启动应用上下文
	if err := ctx.Start(); err != nil {
		fmt.Printf("Failed to start context: %v\n", err)
		return
	}

	// 获取服务并使用
	retrievedOrderService := ctx.GetBean("orderService")
	if retrievedOrderService == nil {
		fmt.Printf("Failed to get order service\n")
	} else {
		if orderSvc, ok := retrievedOrderService.(*OrderService); ok {
			order := orderSvc.CreateOrder(123)
			fmt.Printf("Created order: %s\n", order)
		}
	}

	// 停止应用上下文
	if err := ctx.Stop(); err != nil {
		fmt.Printf("Failed to stop context: %v\n", err)
	}
}