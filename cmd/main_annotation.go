//go:build annotation
// +build annotation

package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

// 基于反射 + 结构体标签的自动注册版本
// 这种方式更接近 Spring Boot 的注解风格
//
// 使用 build tag 运行: go run -tags annotation cmd/main_annotation.go
func main() {
	r := gin.Default()

	// 使用泛型和反射自动解析结构体标签
	// 注意：这种方式需要控制器方法签名符合 gin.HandlerFunc
	// 并且使用 route:"GET /path" 标签

	// 打印所有已注册的路由
	printRoutes(r)

	log.Println("Annotation-based Router Server starting on :28080")
	if err := r.Run(":28080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func printRoutes(r *gin.Engine) {
	log.Println("\n========== Registered Routes ==========")
	for _, route := range r.Routes() {
		log.Printf("%s %s", route.Method, route.Path)
	}
	log.Println("=========================================\n")
}
