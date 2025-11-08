package middleware

import (
	"net/http"
	"seiyuu-chat/router"
	"sync"
	"time"
)

// RateLimiter 简单的速率限制器
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // 每分钟允许的请求数
	cleanup  time.Duration // 清理间隔
}

// Visitor 访客信息
type Visitor struct {
	limiter  *time.Ticker
	lastSeen time.Time
	count    int
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(rate int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		cleanup:  5 * time.Minute,
	}

	// 启动清理goroutine
	go rl.cleanupVisitors()

	return rl
}

// GetVisitor 获取或创建访客
func (rl *RateLimiter) GetVisitor(ip string) *Visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &Visitor{
			limiter:  time.NewTicker(time.Minute / time.Duration(rl.rate)),
			lastSeen: time.Now(),
			count:    0,
		}
		rl.visitors[ip] = v
	}

	v.lastSeen = time.Now()
	return v
}

// cleanupVisitors 定期清理过期的访客
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(rl.cleanup)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				v.limiter.Stop()
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(rate int) router.HandlerFunc {
	limiter := NewRateLimiter(rate)

	return func(c *router.Context) {
		ip := c.ClientIP()
		visitor := limiter.GetVisitor(ip)

		select {
		case <-visitor.limiter.C:
			visitor.count++
			c.Next()
		default:
			c.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"success": false,
				"error":   "请求过于频繁，请稍后再试",
			})
			c.Abort()
		}
	}
}

// APIRateLimitMiddleware API专用速率限制（较宽松）
func APIRateLimitMiddleware() router.HandlerFunc {
	return RateLimitMiddleware(60) // 每分钟60次请求
}

// AdminRateLimitMiddleware 管理员API速率限制（更宽松）
func AdminRateLimitMiddleware() router.HandlerFunc {
	return RateLimitMiddleware(120) // 每分钟120次请求
}
