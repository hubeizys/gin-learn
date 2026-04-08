package controller

import (
	"net/http"

	"gin-learn/internal/router"

	"github.com/gin-gonic/gin"
)

// IndexController 首页控制器
type IndexController struct {
}

// NewIndexController 构造函数
func NewIndexController() *IndexController {
	return &IndexController{}
}

// Index 首页处理
func (i *IndexController) Index(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to Gin Auto Router!")
}

// Health 健康检查
func (i *IndexController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "gin-auto-router",
	})
}

// GetRoutes 实现 AutoRoutable 接口
func (i *IndexController) GetRoutes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{Method: "GET", Path: "/", Handler: i.Index},
		{Method: "GET", Path: "/health", Handler: i.Health},
	}
}
