package main

import (
	"log"

	"gin-learn/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建默认 Gin 引擎
	r := gin.Default()

	// 注册路由
	handler.Setup(r)

	// 启动服务器
	log.Println("Server starting on :28080")
	if err := r.Run(":28080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
