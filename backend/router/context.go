package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Context HTTP上下文，替代gin.Context
type Context struct {
	Request    *http.Request
	Writer     http.ResponseWriter
	params     map[string]string
	Context    context.Context
	handlers   []HandlerFunc
	index      int
	statusCode int
	Errors     []error
}

// HandlerFunc 处理函数类型
type HandlerFunc func(*Context)

// NewContext 创建新的上下文
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:    r,
		Writer:     w,
		Context:    r.Context(),
		params:     make(map[string]string),
		index:      -1,
		statusCode: 200, // 默认状态码
		Errors:     make([]error, 0),
	}
}

// Param 获取路径参数
func (c *Context) Param(key string) string {
	return c.params[key]
}

// SetParam 设置路径参数
func (c *Context) SetParam(key, value string) {
	c.params[key] = value
}

// Query 获取查询参数
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// DefaultQuery 获取查询参数，带默认值
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value := c.Query(key); value != "" {
		return value
	}
	return defaultValue
}

// GetHeader 获取请求头
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// ShouldBindJSON 绑定JSON数据
func (c *Context) ShouldBindJSON(obj interface{}) error {
	if c.Request.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(obj)
}

// JSON 返回JSON响应
func (c *Context) JSON(code int, obj interface{}) {
	c.statusCode = code
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)

	encoder := json.NewEncoder(c.Writer)
	encoder.Encode(obj)
}

// Set 设置上下文值
func (c *Context) Set(key string, value interface{}) {
	c.Context = context.WithValue(c.Context, key, value)
}

// Get 获取上下文值
func (c *Context) Get(key string) (interface{}, bool) {
	value := c.Context.Value(key)
	return value, value != nil
}

// GetString 获取字符串类型的上下文值
func (c *Context) GetString(key string) string {
	if value, exists := c.Get(key); exists {
		if s, ok := value.(string); ok {
			return s
		}
	}
	return ""
}

// Next 执行下一个中间件/处理器
func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
}

// Abort 中断执行链
func (c *Context) Abort() {
	c.index = len(c.handlers)
}

// AbortWithStatus 中断执行并设置状态码
func (c *Context) AbortWithStatus(code int) {
	c.statusCode = code
	c.Writer.WriteHeader(code)
	c.Abort()
}

// IsAborted 检查是否已中断
func (c *Context) IsAborted() bool {
	return c.index >= len(c.handlers)
}

// setHandlers 设置处理器链（内部使用）
func (c *Context) setHandlers(handlers []HandlerFunc) {
	c.handlers = handlers
}

// Method 获取HTTP方法
func (c *Context) Method() string {
	return c.Request.Method
}

// Path 获取请求路径
func (c *Context) Path() string {
	return c.Request.URL.Path
}

// ClientIP 获取客户端IP地址
func (c *Context) ClientIP() string {
	// 检查 X-Forwarded-For 头
	if xff := c.Request.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 检查 X-Real-IP 头
	if xri := c.Request.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 检查 CF-Connecting-IP 头（Cloudflare）
	if cfip := c.Request.Header.Get("CF-Connecting-IP"); cfip != "" {
		return cfip
	}

	// 使用 RemoteAddr
	ip := c.Request.RemoteAddr
	if strings.Contains(ip, ":") {
		if host, _, found := strings.Cut(ip, ":"); found {
			return host
		}
	}
	return ip
}

// Status 获取HTTP状态码
func (c *Context) Status() int {
	return c.statusCode
}

// UserAgent 获取User-Agent
func (c *Context) UserAgent() string {
	return c.Request.UserAgent()
}

// extractPathParams 从URL路径中提取参数
func extractPathParams(pattern, path string) map[string]string {
	params := make(map[string]string)

	// 分割路径
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return params
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			// 路径参数
			paramName := part[1:]
			params[paramName] = pathParts[i]
		}
	}

	return params
}

// matchPath 检查路径是否匹配模式
func matchPath(pattern, path string) bool {
	// 处理精确匹配
	if !strings.Contains(pattern, ":") {
		return pattern == path
	}

	// 处理参数匹配
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if !strings.HasPrefix(part, ":") && part != pathParts[i] {
			return false
		}
	}

	return true
}