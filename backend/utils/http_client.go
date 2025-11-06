package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// HTTPClientConfig HTTP客户端配置
type HTTPClientConfig struct {
	Timeout         time.Duration
	RetryCount      int
	RetryWaitTime   time.Duration
	RetryMaxWaitTime time.Duration
	UserAgent       string
}

// DefaultHTTPConfig 获取默认的HTTP配置
func DefaultHTTPConfig() *HTTPClientConfig {
	// 从环境变量读取配置，如果没有则使用默认值
	timeout := getEnvDuration("HTTP_TIMEOUT", 30*time.Second)
	retryCount := getEnvInt("HTTP_RETRY_COUNT", 2)
	retryWaitTime := getEnvDuration("HTTP_RETRY_WAIT_TIME", 1*time.Second)
	retryMaxWaitTime := getEnvDuration("HTTP_RETRY_MAX_WAIT_TIME", 5*time.Second)
	userAgent := getEnvString("HTTP_USER_AGENT", "NijiChat/1.0")

	return &HTTPClientConfig{
		Timeout:          timeout,
		RetryCount:       retryCount,
		RetryWaitTime:    retryWaitTime,
		RetryMaxWaitTime: retryMaxWaitTime,
		UserAgent:       userAgent,
	}
}

// NewRestyClient 创建配置好的Resty客户端
func NewRestyClient(config *HTTPClientConfig) *resty.Client {
	if config == nil {
		config = DefaultHTTPConfig()
	}

	client := resty.New()
	client.SetTimeout(config.Timeout)
	client.SetHeader("User-Agent", config.UserAgent)
	client.SetRetryCount(config.RetryCount)
	client.SetRetryWaitTime(config.RetryWaitTime)
	client.SetRetryMaxWaitTime(config.RetryMaxWaitTime)

	// 在开发环境下启用调试模式
	if os.Getenv("GO_ENV") == "development" {
		client.SetDebug(true)
	}

	return client
}

// NewRestyClientWithAuth 创建带认证的Resty客户端
func NewRestyClientWithAuth(config *HTTPClientConfig, authToken string) *resty.Client {
	client := NewRestyClient(config)

	if authToken != "" {
		client.SetHeader("Content-Type", "application/json")
		client.SetAuthToken(authToken) // 自动设置Bearer token
	}

	return client
}

// 辅助函数：从环境变量读取time.Duration，如果不存在则使用默认值
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}

// 辅助函数：从环境变量读取int，如果不存在则使用默认值
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// 辅助函数：从环境变量读取string，如果不存在则使用默认值
func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}