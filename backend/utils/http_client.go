//go:build js && wasm
// +build js,wasm

package utils

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/syumai/workers/cloudflare/fetch"
)

// HTTPClientConfig HTTP客户端配置
type HTTPClientConfig struct {
	Timeout          time.Duration
	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
	UserAgent        string
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
		UserAgent:        userAgent,
	}
}

// CloudflareHTTPClient 专用于 Cloudflare Workers 的 HTTP 客户端包装器
type CloudflareHTTPClient struct {
	client *http.Client
}

// NewCloudflareHTTPClient 创建 Cloudflare Workers 兼容的 HTTP 客户端
func NewCloudflareHTTPClient(config *HTTPClientConfig) *CloudflareHTTPClient {
	if config == nil {
		config = DefaultHTTPConfig()
	}

	// 使用 Cloudflare Workers 的 fetch API
	fetchClient := fetch.NewClient()
	httpClient := fetchClient.HTTPClient(fetch.RedirectModeFollow)

	// 设置超时
	httpClient.Timeout = config.Timeout

	return &CloudflareHTTPClient{
		client: httpClient,
	}
}

// NewCloudflareHTTPClientWithAuth 创建带认证的 Cloudflare Workers HTTP 客户端
func NewCloudflareHTTPClientWithAuth(config *HTTPClientConfig, authToken string) *CloudflareHTTPClient {
	baseClient := NewCloudflareHTTPClient(config)

	if authToken != "" {
		// 包装Transport以添加认证头
		baseClient.client.Transport = &authTransport{
			base:      baseClient.client.Transport,
			authToken: authToken,
		}
	}

	return baseClient
}

// HTTPClient 返回标准的 http.Client
func (c *CloudflareHTTPClient) HTTPClient() *http.Client {
	return c.client
}

// authTransport 是一个包装器，用于在每个请求中添加认证头
type authTransport struct {
	base      http.RoundTripper
	authToken string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 克隆请求以避免修改原始请求
	newReq := req.Clone(req.Context())

	// 添加认证头
	newReq.Header.Set("Authorization", "Bearer "+t.authToken)

	return t.base.RoundTrip(newReq)
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

// LogHTTPRequest 记录HTTP请求日志（仅在开发环境）
func LogHTTPRequest(service, operation, method, url string) {
	if os.Getenv("GO_ENV") == "development" {
		log.Printf("[%s:%s] %s %s", service, operation, method, url)
	}
}
