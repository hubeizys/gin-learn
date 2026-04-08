# Gin 自动路由注册方案

Gin 本身没有内置像 Spring Boot 那样的注解驱动路由，但可以通过以下方式实现类似功能：

## 方案一：AutoRoutable 接口方式（推荐）✅

类似 Spring Boot 的 `@RestController` 注解，每个控制器实现 `GetRoutes()` 方法返回路由定义。

### 使用方式

**1. 定义控制器，实现 AutoRoutable 接口：**

```go
type UserController struct{}

func (u *UserController) GetRoutes() []router.RouteDefinition {
    return []router.RouteDefinition{
        {Method: "GET", Path: "/api/users", Handler: u.GetUsers},
        {Method: "GET", Path: "/api/users/:id", Handler: u.GetUser},
        {Method: "POST", Path: "/api/users", Handler: u.CreateUser},
    }
}
```

**2. 注册并自动路由：**

```go
func main() {
    r := gin.Default()
    
    // 实例化控制器并注册
    router.RegisterController(controller.NewUserController())
    router.RegisterController(controller.NewOrderController())
    
    // 一行代码完成所有路由注册
    router.AutoRegister(r)
    
    r.Run(":8080")
}
```

### 优点
- ✅ 类型安全，IDE 友好
- ✅ 编译时检查
- ✅ 路由定义清晰可见
- ✅ 支持动态路由前缀

### 缺点
- ❌ 需要手动定义路由
- ❌ 路由和处理器分离

---

## 方案二：结构体标签 + 反射方式

通过在方法上添加注释标签，然后使用反射自动解析。

### 使用方式

```go
// route:"group:/api"
type UserController struct{}

// route:"GET /users"
func (u *UserController) GetUsers(c *gin.Context) { ... }

// route:"POST /users"
func (u *UserController) CreateUser(c *gin.Context) { ... }
```

### 实现原理

```go
// 基于方法名自动映射 HTTP 方法
switch {
case strings.HasPrefix(methodName, "Get"):
    httpMethod = "GET"
case strings.HasPrefix(methodName, "Post"):
    httpMethod = "POST"
case strings.HasPrefix(methodName, "Put"):
    httpMethod = "PUT"
case strings.HasPrefix(methodName, "Delete"):
    httpMethod = "DELETE"
}

// 路径自动生成：方法名转小写
path = "/" + strings.ToLower(methodName[3:])  // GetUsers -> /users
```

### 优点
- ✅ 更接近 Spring Boot 风格
- ✅ 路由直接在方法上声明
- ✅ 代码更紧凑

### 缺点
- ❌ 依赖反射，性能略低
- ❌ IDE 提示较弱
- ❌ 路径命名需要遵循约定

---

## 方案三：第三方框架

如果你想要更完整的注解支持，可以使用社区框架：

### 1. go-kits/gin
```go
@router.Get("/users")
func (u *UserController) GetUsers() { ... }
```

### 2. swaggo/gin-swagger
```go
// @Summary 获取用户列表
// @Router /api/users [get]
func GetUsers() { ... }
```

### 3. gin-vue-admin 后台管理框架
使用 YAML 配置驱动路由

---

## 项目文件结构

```
gin-learn/
├── cmd/
│   ├── main.go              # 自动路由注册示例
│   ├── main_auto.go         # 完整控制器示例
│   └── main_annotation.go   # 反射注解示例
├── internal/
│   ├── router/
│   │   └── annotations.go   # 路由注册核心逻辑
│   └── controller/
│       ├── user_controller.go
│       ├── order_controller.go
│       └── index_controller.go
└── docs/
    └── auto-router.md
```

---

## 运行方式

```bash
# 默认版本（推荐）
go run ./cmd/main.go

# 完整控制器示例
go run -tags auto ./cmd/main_auto.go

# 反射注解示例
go run -tags annotation ./cmd/main_annotation.go
```

---

## 对比总结

| 特性 | AutoRoutable 接口 | 结构体标签+反射 | 手动注册 |
|------|-------------------|----------------|----------|
| 类型安全 | ✅ | ⚠️ | ✅ |
| IDE 友好 | ✅ | ⚠️ | ✅ |
| 性能 | ✅ | ⚠️ | ✅ |
| 代码量 | 中 | 少 | 多 |
| 维护性 | ✅ | ⚠️ | ✅ |
| Spring Boot 风格 | ⚠️ | ✅ | ❌ |

---

## 建议

**小型项目**：使用 `AutoRoutable` 接口方式，清晰且类型安全

**中型项目**：可以封装一个更完整的自动注册框架

**大型项目**：考虑使用 `go-kits` 或其他成熟的注解框架

如果你习惯 Spring Boot 的注解风格，**方案二** 是最接近的，但建议在团队内部维护统一的代码规范。
