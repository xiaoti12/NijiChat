package router

import (
	"net/http"
)

// Engine 路由引擎
type Engine struct {
	routes         []Route
	middlewares    []HandlerFunc
	noRouteHandler HandlerFunc
}

// Route 路由定义
type Route struct {
	Method   string
	Pattern  string
	Handlers []HandlerFunc
}

// RouterGroup 路由组
type RouterGroup struct {
	engine      *Engine
	prefix      string
	middlewares []HandlerFunc
}

// New 创建新的路由引擎
func New() *Engine {
	return &Engine{
		routes:      make([]Route, 0),
		middlewares: make([]HandlerFunc, 0),
	}
}

// Use 添加全局中间件
func (e *Engine) Use(middleware HandlerFunc) {
	e.middlewares = append(e.middlewares, middleware)
}

// Group 创建路由组
func (e *Engine) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		engine:      e,
		prefix:      prefix,
		middlewares: make([]HandlerFunc, 0),
	}
}

// NoRoute 设置404处理器
func (e *Engine) NoRoute(handler HandlerFunc) {
	e.noRouteHandler = handler
}

// ServeHTTP 实现http.Handler接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	e.handleRequest(c)
}

// handleRequest 处理HTTP请求
func (e *Engine) handleRequest(c *Context) {
	// 查找匹配的路由
	for _, route := range e.routes {
		if route.Method == c.Method() && matchPath(route.Pattern, c.Path()) {
			// 提取路径参数
			params := extractPathParams(route.Pattern, c.Path())
			for key, value := range params {
				c.SetParam(key, value)
			}

			// 合并全局中间件和路由处理器
			handlers := make([]HandlerFunc, 0, len(e.middlewares)+len(route.Handlers))
			handlers = append(handlers, e.middlewares...)
			handlers = append(handlers, route.Handlers...)

			c.setHandlers(handlers)
			c.Next()
			return
		}
	}

	// 特殊处理OPTIONS预检请求：如果存在相同路径的其他方法路由，允许OPTIONS请求通过到中间件处理
	if c.Method() == "OPTIONS" {
		for _, route := range e.routes {
			if route.Method != "OPTIONS" && matchPath(route.Pattern, c.Path()) {
				// 找到匹配的路径，只使用全局中间件处理OPTIONS请求
				// 中间件中的CORS处理器会处理这个OPTIONS请求
				handlers := make([]HandlerFunc, 0, len(e.middlewares))
				handlers = append(handlers, e.middlewares...)

				c.setHandlers(handlers)
				c.Next()
				return
			}
		}
	}

	// 没有找到匹配的路由
	if e.noRouteHandler != nil {
		c.setHandlers([]HandlerFunc{e.noRouteHandler})
		c.Next()
	} else {
		c.Writer.WriteHeader(http.StatusNotFound)
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "接口不存在",
		})
	}
}

// GET 注册GET路由
func (e *Engine) GET(pattern string, handlers ...HandlerFunc) {
	e.addRoute("GET", pattern, handlers)
}

// POST 注册POST路由
func (e *Engine) POST(pattern string, handlers ...HandlerFunc) {
	e.addRoute("POST", pattern, handlers)
}

// PUT 注册PUT路由
func (e *Engine) PUT(pattern string, handlers ...HandlerFunc) {
	e.addRoute("PUT", pattern, handlers)
}

// DELETE 注册DELETE路由
func (e *Engine) DELETE(pattern string, handlers ...HandlerFunc) {
	e.addRoute("DELETE", pattern, handlers)
}

// OPTIONS 注册OPTIONS路由
func (e *Engine) OPTIONS(pattern string, handlers ...HandlerFunc) {
	e.addRoute("OPTIONS", pattern, handlers)
}

// addRoute 添加路由
func (e *Engine) addRoute(method, pattern string, handlers []HandlerFunc) {
	route := Route{
		Method:   method,
		Pattern:  pattern,
		Handlers: handlers,
	}
	e.routes = append(e.routes, route)
}

// Use 添加组中间件
func (rg *RouterGroup) Use(middleware HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middleware)
}

// Group 创建子路由组
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		engine:      rg.engine,
		prefix:      rg.prefix + prefix,
		middlewares: append([]HandlerFunc{}, rg.middlewares...), // 继承父组的中间件
	}
}

// GET 组GET路由
func (rg *RouterGroup) GET(pattern string, handlers ...HandlerFunc) {
	rg.addRoute("GET", pattern, handlers)
}

// POST 组POST路由
func (rg *RouterGroup) POST(pattern string, handlers ...HandlerFunc) {
	rg.addRoute("POST", pattern, handlers)
}

// PUT 组PUT路由
func (rg *RouterGroup) PUT(pattern string, handlers ...HandlerFunc) {
	rg.addRoute("PUT", pattern, handlers)
}

// DELETE 组DELETE路由
func (rg *RouterGroup) DELETE(pattern string, handlers ...HandlerFunc) {
	rg.addRoute("DELETE", pattern, handlers)
}

// OPTIONS 组OPTIONS路由
func (rg *RouterGroup) OPTIONS(pattern string, handlers ...HandlerFunc) {
	rg.addRoute("OPTIONS", pattern, handlers)
}

// addRoute 添加组路由
func (rg *RouterGroup) addRoute(method, pattern string, handlers []HandlerFunc) {
	fullPattern := rg.prefix + pattern

	// 合并组中间件和处理器
	allHandlers := make([]HandlerFunc, 0, len(rg.middlewares)+len(handlers))
	allHandlers = append(allHandlers, rg.middlewares...)
	allHandlers = append(allHandlers, handlers...)

	rg.engine.addRoute(method, fullPattern, allHandlers)
}
