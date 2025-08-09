package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"gospring/context"
)

// 定义服务接口
type ProductService interface {
	GetProduct(id int) *Product
	GetAllProducts() []*Product
	CreateProduct(name string, price float64) *Product
}

type OrderService interface {
	CreateOrder(productId int, quantity int) *Order
	GetOrder(id int) *Order
}

// 模型定义
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Order struct {
	ID        int      `json:"id"`
	ProductID int      `json:"product_id"`
	Quantity  int      `json:"quantity"`
	Total     float64  `json:"total"`
	Product   *Product `json:"product"`
}

// ProductServiceImpl 产品服务实现
type ProductServiceImpl struct {
	products map[int]*Product
	nextID   int
	
	_ string `component:"productService" singleton:"true"`
}

func (s *ProductServiceImpl) GetProduct(id int) *Product {
	return s.products[id]
}

func (s *ProductServiceImpl) GetAllProducts() []*Product {
	var products []*Product
	for _, product := range s.products {
		products = append(products, product)
	}
	return products
}

func (s *ProductServiceImpl) CreateProduct(name string, price float64) *Product {
	s.nextID++
	product := &Product{
		ID:    s.nextID,
		Name:  name,
		Price: price,
	}
	s.products[product.ID] = product
	return product
}

func (s *ProductServiceImpl) Init() error {
	s.products = make(map[int]*Product)
	s.nextID = 0
	
	// 初始化一些测试数据
	s.CreateProduct("笔记本电脑", 5999.99)
	s.CreateProduct("无线鼠标", 199.99)
	s.CreateProduct("机械键盘", 599.99)
	
	fmt.Println("ProductService 初始化完成")
	return nil
}

// OrderServiceImpl 订单服务实现
type OrderServiceImpl struct {
	ProductService ProductService `inject:"productService"`
	orders         map[int]*Order
	nextID         int
	
	_ string `component:"orderService" singleton:"true"`
}

func (s *OrderServiceImpl) CreateOrder(productId int, quantity int) *Order {
	product := s.ProductService.GetProduct(productId)
	if product == nil {
		return nil
	}
	
	s.nextID++
	order := &Order{
		ID:        s.nextID,
		ProductID: productId,
		Quantity:  quantity,
		Total:     product.Price * float64(quantity),
		Product:   product,
	}
	
	s.orders[order.ID] = order
	return order
}

func (s *OrderServiceImpl) GetOrder(id int) *Order {
	return s.orders[id]
}

func (s *OrderServiceImpl) Init() error {
	s.orders = make(map[int]*Order)
	s.nextID = 0
	fmt.Println("OrderService 初始化完成")
	return nil
}

// ProductController 产品控制器
type ProductController struct {
	ProductService ProductService `inject:"productService"`
	
	_ string `component:"productController" singleton:"true"`
}

func (c *ProductController) SetupRoutes() {
	http.HandleFunc("/products", c.handleProducts)
	http.HandleFunc("/products/", c.handleProduct)
}

func (c *ProductController) handleProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "GET":
		products := c.ProductService.GetAllProducts()
		json.NewEncoder(w).Encode(products)
		
	case "POST":
		var req struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		product := c.ProductService.CreateProduct(req.Name, req.Price)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *ProductController) handleProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 从URL路径中提取ID
	path := r.URL.Path
	idStr := path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	
	switch r.Method {
	case "GET":
		product := c.ProductService.GetProduct(id)
		if product == nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(product)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// OrderController 订单控制器
type OrderController struct {
	OrderService OrderService `inject:"orderService"`
	
	_ string `component:"orderController" singleton:"true"`
}

func (c *OrderController) SetupRoutes() {
	http.HandleFunc("/orders", c.handleOrders)
	http.HandleFunc("/orders/", c.handleOrder)
}

func (c *OrderController) handleOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "POST":
		var req struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		order := c.OrderService.CreateOrder(req.ProductID, req.Quantity)
		if order == nil {
			http.Error(w, "Product not found", http.StatusBadRequest)
			return
		}
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *OrderController) handleOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 从URL路径中提取ID
	path := r.URL.Path
	idStr := path[len("/orders/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	
	switch r.Method {
	case "GET":
		order := c.OrderService.GetOrder(id)
		if order == nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(order)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	fmt.Println("=== GoSpring Web 应用演示 ===")
	
	// 创建应用上下文
	ctx := context.NewApplicationContext()
	
	// 创建服务实例
	productService := &ProductServiceImpl{}
	orderService := &OrderServiceImpl{}
	
	// 创建控制器实例
	productController := &ProductController{}
	orderController := &OrderController{}
	
	// 注册组件
	fmt.Println("1. 注册组件...")
	if err := ctx.RegisterComponents(
		productService, 
		orderService, 
		productController, 
		orderController,
	); err != nil {
		log.Fatalf("注册组件失败: %v", err)
	}
	
	// 通过接口注册服务
	productServiceType := reflect.TypeOf((*ProductService)(nil)).Elem()
	orderServiceType := reflect.TypeOf((*OrderService)(nil)).Elem()
	
	ctx.RegisterByInterface(productServiceType, productService, "productServiceInterface")
	ctx.RegisterByInterface(orderServiceType, orderService, "orderServiceInterface")
	
	// 启动上下文
	fmt.Println("2. 启动应用上下文...")
	if err := ctx.Start(); err != nil {
		log.Fatalf("启动上下文失败: %v", err)
	}
	
	// 设置路由
	fmt.Println("3. 设置HTTP路由...")
	
	// 获取控制器并设置路由
	prodCtrl := ctx.GetBean("productController").(*ProductController)
	orderCtrl := ctx.GetBean("orderController").(*OrderController)
	
	prodCtrl.SetupRoutes()
	orderCtrl.SetupRoutes()
	
	// 添加根路径处理
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"message": "欢迎使用 GoSpring Web API",
				"endpoints": map[string]string{
					"GET /products":     "获取所有产品",
					"POST /products":    "创建新产品",
					"GET /products/{id}": "获取指定产品",
					"POST /orders":      "创建新订单",
					"GET /orders/{id}":  "获取指定订单",
				},
				"example_requests": map[string]interface{}{
					"create_product": map[string]interface{}{
						"url":    "POST /products",
						"body":   map[string]interface{}{"name": "新产品", "price": 99.99},
					},
					"create_order": map[string]interface{}{
						"url":    "POST /orders",
						"body":   map[string]interface{}{"product_id": 1, "quantity": 2},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			http.NotFound(w, r)
		}
	})
	
	// 显示注册的Bean信息
	fmt.Println("\n4. 注册的组件:")
	beans := ctx.ListBeans()
	for _, beanName := range beans {
		beanDef := ctx.GetBeanDefinition(beanName)
		fmt.Printf("  - %s (类型: %v)\n", beanName, beanDef.Type)
	}
	
	// 启动HTTP服务器
	port := ":8080"
	fmt.Printf("\n5. 启动HTTP服务器，监听端口 %s\n", port)
	fmt.Println("API端点:")
	fmt.Println("  GET  http://localhost:8080/          - API文档")
	fmt.Println("  GET  http://localhost:8080/products  - 获取所有产品")
	fmt.Println("  POST http://localhost:8080/products  - 创建产品")
	fmt.Println("  GET  http://localhost:8080/products/1 - 获取产品1")
	fmt.Println("  POST http://localhost:8080/orders    - 创建订单")
	fmt.Println("  GET  http://localhost:8080/orders/1  - 获取订单1")
	fmt.Println("\n按 Ctrl+C 停止服务器")
	
	// 启动服务器
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("启动HTTP服务器失败: %v", err)
	}
}