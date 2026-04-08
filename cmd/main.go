package main

import (
	"log"

	router "gin-learn/internal/router"

	"github.com/gin-gonic/gin"
)

// Gin 自动路由注册 - 多种实现方式对比
//
// 主要有两种自动路由注册方案：
//
// 1. AutoRoutable 接口方式（推荐）
//    - 实现 GetRoutes() 方法返回路由定义
//    - 优点：类型安全 IDE 友好
//    - 缺点：需要手动定义路由
//
// 2. 结构体标签 + 反射方式
//    - 在方法上添加 route:"GET /path" 标签
//    - 优点：更接近 Spring Boot 风格
//    - 缺点：依赖反射，IDE 提示较弱
//
// 使用 main_auto.go 查看接口方式
// 使用 main_annotation.go 查看标签方式
func main() {
	r := gin.Default()

	// 自动路由注册示例
	router.RegisterController(&ExampleController{
		BasePath: "/api",
	})

	router.AutoRegister(r)

	// 打印所有已注册的路由
	printRegisteredRoutes(r)

	log.Println("Server starting on :28080")
	if err := r.Run(":28080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// ExampleController 示例控制器
type ExampleController struct {
	BasePath string
}

func (e *ExampleController) GetRoutes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{Method: "GET", Path: e.BasePath + "/hello", Handler: e.Hello},
		{Method: "GET", Path: e.BasePath + "/users/:id", Handler: e.GetUser},
		{Method: "POST", Path: e.BasePath + "/users", Handler: e.CreateUser},
	}
}

func (e *ExampleController) Hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello from auto router!"})
}

func (e *ExampleController) GetUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{"id": id, "name": "User " + id})
}

func (e *ExampleController) CreateUser(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "User created", "name": req.Name})
}

func printRegisteredRoutes(r *gin.Engine) {
	log.Println("\n========== Auto-Registered Routes ==========")
	for _, route := range r.Routes() {
		log.Printf("  %s %s", route.Method, route.Path)
	}
	log.Println("==============================================\n")
}
