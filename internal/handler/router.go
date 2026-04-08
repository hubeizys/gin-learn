package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup 注册所有路由
func Setup(r *gin.Engine) {
	// 首页
	r.GET("/", indexHandler)

	// API 路由组
	api := r.Group("/api")
	{
		api.GET("/hello", helloHandler)
		api.GET("/users/:id", getUserHandler)
		api.POST("/users", createUserHandler)
		api.GET("/query", queryHandler)
	}

	// JSON API 示例
	r.GET("/json", jsonHandler)

	// 绑定 JSON body 示例
	r.POST("/login", loginHandler)
}

// indexHandler 首页处理函数
func indexHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to Gin!")
}

// helloHandler 返回简单的 JSON 响应
func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, Gin!",
	})
}

// User 结构体示例
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// getUserHandler 路径参数示例
func getUserHandler(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": "John Doe",
	})
}

// createUserHandler POST 请求示例
func createUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"user":    user,
	})
}

// queryHandler 查询参数示例
func queryHandler(c *gin.Context) {
	name := c.Query("name")
	age := c.Query("age")
	c.JSON(http.StatusOK, gin.H{
		"name": name,
		"age":  age,
	})
}

// jsonHandler 返回结构体 JSON
func jsonHandler(c *gin.Context) {
	user := User{
		ID:    "1",
		Name:  "Alice",
		Email: "alice@example.com",
	}
	c.JSON(http.StatusOK, user)
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// loginHandler 绑定 JSON body 示例
func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"username": req.Username,
	})
}
