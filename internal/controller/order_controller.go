package controller

import (
	"net/http"

	"gin-learn/internal/router"

	"github.com/gin-gonic/gin"
)

// OrderController 订单控制器
type OrderController struct {
}

// NewOrderController 构造函数
func NewOrderController() *OrderController {
	return &OrderController{}
}

// GetOrders 获取订单列表
func (o *OrderController) GetOrders(c *gin.Context) {
	orders := []gin.H{
		{"id": "1", "product": "Apple", "price": 100},
		{"id": "2", "product": "Banana", "price": 50},
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// GetOrder 获取单个订单
func (o *OrderController) GetOrder(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"product": "Sample Product",
		"price":   100,
	})
}

// CreateOrder 创建订单
func (o *OrderController) CreateOrder(c *gin.Context) {
	var req struct {
		Product string `json:"product" binding:"required"`
		Price   int    `json:"price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created",
		"order":   req,
	})
}

// GetRoutes 实现 AutoRoutable 接口
// route:"group:/api"
func (o *OrderController) GetRoutes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{Method: "GET", Path: "/api/orders", Handler: o.GetOrders},
		{Method: "GET", Path: "/api/orders/:id", Handler: o.GetOrder},
		{Method: "POST", Path: "/api/orders", Handler: o.CreateOrder},
	}
}
