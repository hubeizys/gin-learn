//go:build auto
// +build auto

package main

import (
	"log"

	"gin-learn/internal/controller"
	router "gin-learn/internal/router"

	"github.com/gin-gonic/gin"
)

// 自动路由注册版本 - 类似 Spring Boot 风格
//
// 使用 build tag 运行: go run -tags auto cmd/main_auto.go
func main() {
	// 创建默认 Gin 引擎
	r := gin.Default()

	// ============ 自动路由注册 ============
	// 只需要实例化控制器，就会自动注册路由
	router.RegisterController(controller.NewIndexController())
	router.RegisterController(controller.NewUserController())
	router.RegisterController(controller.NewOrderController())

	// 一行代码完成所有路由注册
	router.AutoRegister(r)

	// 打印所有已注册的路由（方便调试）
	printRoutes(r)

	// 启动服务器
	log.Println("Auto Router Server starting on :28080")
	if err := r.Run(":28080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// printRoutes 打印所有已注册的路由
func printRoutes(r *gin.Engine) {
	log.Println("\n========== Registered Routes ==========")
	for _, route := range r.Routes() {
		log.Printf("%s %s", route.Method, route.Path)
	}
	log.Println("=========================================\n")
}
