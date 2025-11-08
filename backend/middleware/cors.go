package middleware

import (
	"seiyuu-chat/router"
)

// CORSMiddleware CORS中间件
func CORSMiddleware() router.HandlerFunc {
	return func(c *router.Context) {
		// 允许所有来源（生产环境应该配置具体的域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// 允许浏览器访问的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		// 预检请求的缓存时间（秒）
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")

		// 允许携带认证信息
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ConfigurableCORSMiddleware 可配置的CORS中间件
func ConfigurableCORSMiddleware(allowedOrigins []string) router.HandlerFunc {
	return func(c *router.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查来源是否在允许列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				allowed = true
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		if !allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins[0])
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "3600")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
