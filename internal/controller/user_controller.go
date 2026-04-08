package controller

import (
	"net/http"

	"gin-learn/internal/router"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器 - 使用注解自动注册路由
type UserController struct {
}

// NewUserController 构造函数
func NewUserController() *UserController {
	return &UserController{}
}

// GetUsers 获取用户列表
// route:"GET /users"
func (u *UserController) GetUsers(c *gin.Context) {
	users := []gin.H{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUser 获取单个用户
// route:"GET /users/:id"
func (u *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": "User " + id,
	})
}

// CreateUser 创建用户
// route:"POST /users"
func (u *UserController) CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"user":    req,
	})
}

// UpdateUser 更新用户
// route:"PUT /users/:id"
func (u *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated",
		"id":      id,
		"user":    req,
	})
}

// DeleteUser 删除用户
// route:"DELETE /users/:id"
func (u *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted",
		"id":      id,
	})
}

// GetRoutes 实现 AutoRoutable 接口
// route:"group:/api"
func (u *UserController) GetRoutes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{Method: "GET", Path: "/api/users", Handler: u.GetUsers},
		{Method: "GET", Path: "/api/users/:id", Handler: u.GetUser},
		{Method: "POST", Path: "/api/users", Handler: u.CreateUser},
		{Method: "PUT", Path: "/api/users/:id", Handler: u.UpdateUser},
		{Method: "DELETE", Path: "/api/users/:id", Handler: u.DeleteUser},
	}
}
