package router

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// RouteDefinition 路由定义
type RouteDefinition struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

// GetRoutes 获取所有路由（供自动注册使用）
// 每个控制器实现这个接口即可自动注册
type AutoRoutable interface {
	GetRoutes() []RouteDefinition
}

// 存储所有实现了 AutoRoutable 接口的控制器实例
var controllers = make([]AutoRoutable, 0)

// RegisterController 注册控制器
func RegisterController(c AutoRoutable) {
	controllers = append(controllers, c)
}

// AutoRegister 自动注册所有控制器路由
func AutoRegister(r *gin.Engine) {
	for _, ctrl := range controllers {
		routes := ctrl.GetRoutes()
		for _, route := range routes {
			r.Handle(route.Method, route.Path, route.Handler)
		}
	}
}

// ============ 基于反射的注解方式 ============

// MethodEnum HTTP 方法枚举
type MethodEnum string

const (
	MethodGet    MethodEnum = "GET"
	MethodPost   MethodEnum = "POST"
	MethodPut    MethodEnum = "PUT"
	MethodDelete MethodEnum = "DELETE"
	MethodPatch  MethodEnum = "PATCH"
	MethodAny    MethodEnum = "ANY"
)

// ControllerInfo 控制器信息（用于反射解析）
type ControllerInfo struct {
	Prefix string
	Routes map[string]struct {
		Method  string
		Path    string
		Handler string // 函数名
	}
}

// Group 指定路由组前缀
type Group string

// RegisterFromStruct 使用结构体标签自动注册路由
// 要求：控制器方法签名必须符合 gin.HandlerFunc
// 标签格式: route:"GET /path"
func RegisterFromStruct[R any](r *gin.Engine, basePath string) {
	var instance R
	t := reflect.TypeOf(instance)

	// 如果是指针类型，获取元素类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 获取 Group 前缀
	groupPrefix := basePath
	if gp, ok := getGroupTag(t); ok {
		groupPrefix = basePath + gp
	}

	// 遍历所有方法
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if handler, path, httpMethod, ok := parseMethodTag(method); ok {
			fullPath := groupPrefix + path
			r.Handle(httpMethod, fullPath, handler)
		}
	}
}

func getGroupTag(t reflect.Type) (string, bool) {
	// 检查结构体是否有 route 标签
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag, ok := field.Tag.Lookup("route"); ok {
			parts := strings.Split(tag, ";")
			for _, part := range parts {
				if strings.HasPrefix(part, "group:") {
					return strings.TrimPrefix(part, "group:"), true
				}
			}
		}
	}
	return "", false
}

func parseMethodTag(method reflect.Method) (gin.HandlerFunc, string, string, bool) {
	// 获取方法的标签 - 通过获取函数类型来访问 tag
	fnType := method.Type
	
	// 检查是否是有效的方法签名 (func(*Controller, *gin.Context)
	if fnType.NumIn() != 2 {
		return nil, "", "", false
	}
	
	// 第一个参数应该是 *gin.Context
	ginContextType := reflect.TypeOf((*gin.Context)(nil)).Elem()
	if fnType.In(1) != ginContextType {
		return nil, "", "", false
	}
	
	// 检查函数是否是 gin.HandlerFunc
	if fnType.NumOut() != 1 || fnType.Out(0) != reflect.TypeOf((*interface{})(nil)).Elem() {
		return nil, "", "", false
	}

	// 尝试从文档字符串中解析路由信息
	methodName := method.Name
	
	// 简单映射：方法名以 HTTP 方法开头
	var httpMethod, path string
	
	switch {
	case strings.HasPrefix(methodName, "Get"):
		httpMethod = "GET"
		path = "/" + strings.ToLower(methodName[3:])
	case strings.HasPrefix(methodName, "Post"):
		httpMethod = "POST"
		path = "/" + strings.ToLower(methodName[4:])
	case strings.HasPrefix(methodName, "Put"):
		httpMethod = "PUT"
		path = "/" + strings.ToLower(methodName[3:])
	case strings.HasPrefix(methodName, "Delete"):
		httpMethod = "DELETE"
		path = "/" + strings.ToLower(methodName[6:])
	case strings.HasPrefix(methodName, "Patch"):
		httpMethod = "PATCH"
		path = "/" + strings.ToLower(methodName[5:])
	default:
		// 不支持的方法
		return nil, "", "", false
	}
	
	// 获取函数值
	fn := method.Func.Interface()
	if handler, ok := fn.(gin.HandlerFunc); ok {
		return handler, path, httpMethod, true
	}

	return nil, "", "", false
}
