# Gin 工作原理详解

## 目录
- [没有 Gin：原生 net/http](#一没有-gin原生-nethttp)
- [有 Gin：简洁优雅](#二有-gin简洁优雅)
- [Gin 是如何工作的](#三gin-是如何工作的)
- [核心概念对比](#四核心概念对比)

---

## 一、没有 Gin：原生 net/http

### 1.1 最简单的服务器

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // 注册路由和处理函数
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello World")
    })

    // 启动服务器
    http.ListenAndServe(":8080", nil)
}
```

### 1.2 路由和参数处理

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // 手动解析路由模式
    http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"message": "Hello"})
    })

    // 路径参数需要自己解析
    http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
        // /api/users/123 -> 需要手动提取 "123"
        path := r.URL.Path
        id := path[len("/api/users/"):]

        user := User{ID: id, Name: "John", Email: "john@example.com"}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    })

    // 查询参数需要自己解析
    http.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
        // ?name=Alice&age=25 -> 需要手动提取
        name := r.URL.Query().Get("name")
        age := r.URL.Query().Get("age")

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "name": name,
            "age":  age,
        })
    })

    http.ListenAndServe(":8080", nil)
}
```

### 1.3 POST 请求处理

```go
package main

import (
    "encoding/json"
    "io"
    "net/http"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func main() {
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // 设置正确的 Content-Type 检查
        if r.Header.Get("Content-Type") != "application/json" {
            http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
            return
        }

        // 手动读取 body
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }
        defer r.Body.Close()

        // 手动解析 JSON
        var req LoginRequest
        if err := json.Unmarshal(body, &req); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }

        // 手动验证
        if req.Username == "" || req.Password == "" {
            http.Error(w, "username and password required", http.StatusBadRequest)
            return
        }

        // 手动序列化响应
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "message":  "Login successful",
            "username": req.Username,
        })
    })

    http.ListenAndServe(":8080", nil)
}
```

### 1.4 问题总结

```
❌ 路由匹配简陋 (不支持路径参数语法)
❌ 路径参数需要手动字符串处理
❌ 查询参数需要手动解析
❌ JSON 绑定需要手动 io.ReadAll + json.Unmarshal
❌ 验证逻辑手写
❌ 中间件支持有限
❌ 没有统一的错误处理
❌ 没有日志、恢复等基础功能
```

---

## 二、有 Gin：简洁优雅

### 2.1 路由和参数处理

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    r := gin.Default()  // 包含日志和恢复中间件

    r.GET("/api/hello", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Hello"})
    })

    // 路径参数自动解析
    r.GET("/api/users/:id", func(c *gin.Context) {
        id := c.Param("id")  // 直接获取路径参数
        c.JSON(http.StatusOK, User{
            ID:   id,
            Name: "John",
            Email: "john@example.com",
        })
    })

    // 查询参数一行搞定
    r.GET("/api/query", func(c *gin.Context) {
        name := c.Query("name")  // 自动解析 ?name=xxx
        age := c.Query("age")     // 自动解析 ?age=xxx
        c.JSON(http.StatusOK, gin.H{"name": name, "age": age})
    })

    r.Run(":8080")
}
```

### 2.2 POST 请求处理

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func main() {
    r := gin.Default()

    r.POST("/login", func(c *gin.Context) {
        var req LoginRequest
        // 自动读取 body、自动解析 JSON、自动验证 required 标签
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message":  "Login successful",
            "username": req.Username,
        })
    })

    r.Run(":8080")
}
```

### 2.3 中间件

```go
package main

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()  // 处理请求
        latency := time.Since(start)
        println("Request:", c.Request.URL.Path, "Duration:", latency)
    }
}

func main() {
    r := gin.New()

    // 全局中间件
    r.Use(gin.Recovery())  // Panic 恢复
    r.Use(Logger())         // 自定义日志

    // 路由组中间件
    api := r.Group("/api")
    api.Use(func(c *gin.Context) {
        // 认证检查
        c.Next()
    })
    {
        api.GET("/hello", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"message": "Hello"})
        })
    }

    r.Run(":8080")
}
```

---

## 三、Gin 是如何工作的

### 3.1 核心组件

```
┌─────────────────────────────────────────────────────────┐
│                      Gin.Engine                         │
│  ┌─────────────────────────────────────────────────┐  │
│  │                  RouterGroup                      │  │
│  │  ┌─────────────────────────────────────────────┐│  │
│  │  │              methodTrees []methodTree        ││  │
│  │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐        ││  │
│  │  │  │   GET   │ │   POST  │ │   PUT   │  ...   ││  │
│  │  │  │  /      │ │  /users │ │  /:id   │        ││  │
│  │  │  │  /api/* │ │  /login │ │         │        ││  │
│  │  │  └────┬────┘ └────┬────┘ └────┬────┘        ││  │
│  │  └───────┼───────────┼───────────┼─────────────┘│  │
│  └──────────┼───────────┼───────────┼───────────────┘  │
│             │           │           │                  │
│             ▼           ▼           ▼                  │
│       ┌─────────┐ ┌─────────┐ ┌─────────┐               │
│       │Handlers│ │Handlers│ │Handlers│               │
│       └─────────┘ └─────────┘ └─────────┘               │
│                                                          │
│  ┌─────────────────────────────────────────────────┐   │
│  │              middlewares []HandlerFunc            │   │
│  │  [Recovery, Logger, CustomAuth, ...]              │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### 3.2 请求处理流程

```
Client Request
      │
      ▼
┌─────────────────────────────────┐
│         gin.Context             │
│  ┌───────────────────────────┐  │
│  │ Request (Method, URL,     │  │
│  │        Header, Body)      │  │
│  │ ResponseWriter            │  │
│  │ Param, Query, Form        │  │
│  │ Keys (middleware data)    │  │
│  └───────────────────────────┘  │
└─────────────────────────────────┘
      │
      ▼
┌─────────────────────────────────┐
│        Middlewares (链式)       │
│  Recovery → Logger → Auth → ... │
│       │         │       │       │
│       └─────────┴───────┘       │
│                 │                │
│                 ▼                │
│      ┌─────────────────────┐    │
│      │     Handler         │    │
│      │ c.JSON() / c.String │    │
│      └─────────────────────┘    │
└─────────────────────────────────┘
      │
      ▼
   Response
```

### 3.3 路由匹配原理

```go
// Gin 内部使用 httprouter (或类似的 Radix Tree)
// 路由 /api/users/:id 的匹配过程：

// 1. 注册路由时构建前缀树
//    /
│── api
    └── users
        └── :id  (参数节点)

// 2. 请求 /api/users/123 匹配过程
//    从根节点开始
//    ├── "/"  → 匹配
//    ├── "api" → 匹配
//    ├── "users" → 匹配
//    └── ":id" → 匹配且提取 "123" 到 Param("id")

// 3. 匹配结果
//    c.Param("id") → "123"
```

### 3.4 中间件链式调用

```go
// 伪代码：Gin 中间件如何工作

func Next(c *Context) {
    // index 记录当前执行到哪个中间件
    c.index++

    if c.index < len(c.middlewares) {
        // 执行下一个中间件
        c.middlewares[c.index](c)
    } else {
        // 所有中间件执行完毕，执行最终处理函数
        c.Handler()
    }
}

// 实际执行顺序：
// middleware[0](c)     // Recovery
//   └── middleware[1](c)   // Logger
//         └── middleware[2](c)   // Auth
//               └── handler()        // 你的业务逻辑
```

---

## 四、核心概念对比

### 4.1 功能对比表

| 功能 | 原生 net/http | Gin |
|------|--------------|-----|
| 路由注册 | `http.HandleFunc` | `r.GET()`, `r.POST()` |
| 路径参数 | 手动字符串处理 | `c.Param("id")` |
| 查询参数 | `r.URL.Query().Get()` | `c.Query("name")` |
| Header | `r.Header.Get()` | `c.GetHeader()` |
| Body 读取 | `io.ReadAll(r.Body)` | `c.ShouldBindJSON()` |
| JSON 响应 | 手动 `json.NewEncoder` | `c.JSON()` |
| 中间件 | 有限支持 | 完整链式中间件 |
| 路由组 | 不支持 | `r.Group()` |
| 参数验证 | 手动 | `binding:"required"` 标签 |
| 错误恢复 | 手动 | 内置 Panic 恢复 |
| 日志 | 手动 | 内置/可自定义 |

### 4.2 代码量对比

| 场景 | 原生代码行数 | Gin 代码行数 |
|------|-------------|-------------|
| GET + 返回 JSON | ~15 行 | ~5 行 |
| POST + 参数验证 | ~40 行 | ~15 行 |
| 中间件链 | ~50 行 | ~10 行 |

### 4.3 选择建议

```
┌─────────────────────────────────────────┐
│  什么时候用原生 net/http？               │
│                                         │
│  ✅ 学习 HTTP 协议底层原理               │
│  ✅ 极简单的场景 (1-2个路由)             │
│  ✅ 不能加外部依赖                       │
│  ✅ 对二进制大小敏感 (嵌入式)            │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│  什么时候用 Gin？                        │
│                                         │
│  ✅ 快速开发 HTTP API                    │
│  ✅ 需要路由参数、查询参数               │
│  ✅ 需要中间件 (认证、日志、限流)         │
│  ✅ 需要参数验证                         │
│  ✅ 需要路由组管理 API 版本              │
│  ✅ 生产级 Web 应用                      │
└─────────────────────────────────────────┘
```

---

## 五、Gin 常用方法速查

```go
r := gin.Default()  // 默认引擎 (带日志+恢复)

// ===== 路由 =====
r.GET("/path", handler)
r.POST("/path", handler)
r.PUT("/path", handler)
r.DELETE("/path", handler)
r.PATCH("/path", handler)
r.OPTIONS("/path", handler)
r.HEAD("/path", handler)
r.Any("/path", handler)           // 匹配任意方法
r.NoRoute(handler)                // 404 处理

// ===== 路由组 =====
v1 := r.Group("/v1")
v1.Use(authMiddleware())          // 路由组中间件
v1.GET("/users", handler)

// ===== 参数获取 =====
c.Param("id")                     // 路径参数 /users/:id
c.Query("name")                   // 查询参数 ?name=xxx
c.PostForm("name")                // 表单参数
c.Header("X-Request-ID")          // 请求头
c.Cookie("session_id")            // Cookie

// ===== 绑定与验证 =====
c.ShouldBindJSON(&struct{})       // JSON body
c.ShouldBindQuery(&struct{})      // Query 参数绑定
c.ShouldBindUri(&struct{})        // 路径参数绑定

// ===== 响应 =====
c.JSON(200, gin.H{"msg": "ok"})
c.String(200, "Hello %s", name)
c.HTML(200, "index.html", data)
c.XML(200, data)
c.File("./static/img.png")
c.Data(200, "application/octet-stream", data)
c.JSONPretty(200, data, "  ")

// ===== 重定向 =====
c.Redirect(301, "https://example.com")
c.Request.URL.Path = "/new-path"
c.Next()  // 内部重定向

// ===== 中间件 =====
r.Use(middleware1)
r.Use(middleware2)
r.GET("/path", middleware3, handler)  // 单路由中间件

// ===== 上下文传递 =====
c.Set("user_id", 123)
c.Get("user_id")  // 在后续 handler/middleware 获取

// ===== 响应处理 =====
c.Next()    // 执行下一个 handler
c.Abort()   // 停止执行后续 handler
c.AbortWithStatus(401)  // 直接返回状态码
```
