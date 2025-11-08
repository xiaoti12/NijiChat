package middleware

import (
	"fmt"
	"seiyuu-chat/router"
	"time"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware() router.HandlerFunc {
	return func(c *router.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqURI := c.Request.RequestURI

		// 状态码
		statusCode := c.Status()

		// 客户端IP
		clientIP := c.ClientIP()

		// 日志输出
		fmt.Printf("[%s] %s | %3d | %13v | %15s | %s\n",
			endTime.Format("2006-01-02 15:04:05"),
			reqMethod,
			statusCode,
			latencyTime,
			clientIP,
			reqURI,
		)
	}
}

// DetailedLoggingMiddleware 详细日志中间件（包含请求体和响应体）
func DetailedLoggingMiddleware() router.HandlerFunc {
	return func(c *router.Context) {
		startTime := time.Now()

		// 打印请求信息
		fmt.Printf("\n[REQUEST] %s %s\n", c.Request.Method, c.Request.RequestURI)
		fmt.Printf("  Client IP: %s\n", c.ClientIP())
		fmt.Printf("  User-Agent: %s\n", c.UserAgent())

		// 打印请求头（可选）
		if c.Request.Header.Get("Authorization") != "" {
			fmt.Printf("  Authorization: [PRESENT]\n")
		}

		c.Next()

		// 计算执行时间
		latency := time.Since(startTime)
		statusCode := c.Status()

		// 打印响应信息
		fmt.Printf("[RESPONSE] Status: %d | Latency: %v\n", statusCode, latency)

		// 打印错误信息（如果有）
		if len(c.Errors) > 0 {
			fmt.Printf("[ERRORS] %v\n", c.Errors)
		}

		fmt.Println()
	}
}

// ErrorLoggingMiddleware 错误日志中间件
func ErrorLoggingMiddleware() router.HandlerFunc {
	return func(c *router.Context) {
		c.Next()

		// 如果有错误，记录详细信息
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				fmt.Printf("[ERROR] %s | %s | %s\n",
					time.Now().Format("2006-01-02 15:04:05"),
					c.Request.RequestURI,
					err.Error(),
				)
			}
		}
	}
}
