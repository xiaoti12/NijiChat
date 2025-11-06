package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

// HTTPError HTTP请求错误类型
type HTTPError struct {
	Service    string // 服务名称（如 "moegirl", "ai"）
	Operation  string // 操作名称（如 "GetRawData", "CallLightweightAI"）
	StatusCode int    // HTTP状态码
	Message    string // 错误消息
	RequestURL string // 请求URL
	Original   error  // 原始错误
}

// Error 实现error接口
func (e *HTTPError) Error() string {
	return fmt.Sprintf("[%s:%s] HTTP %d: %s (URL: %s)",
		e.Service, e.Operation, e.StatusCode, e.Message, e.RequestURL)
}

// Unwrap 支持错误包装
func (e *HTTPError) Unwrap() error {
	return e.Original
}

// NewHTTPError 创建HTTP错误
func NewHTTPError(service, operation string, resp *resty.Response, original error) *HTTPError {
	var message string
	var statusCode int
	var requestURL string

	if resp != nil {
		statusCode = resp.StatusCode()
		requestURL = resp.Request.URL

		if resp.IsError() {
			message = fmt.Sprintf("HTTP error: %s", string(resp.Body()))
		} else {
			message = "Unknown HTTP error"
		}
	} else {
		message = "Network error"
	}

	if original != nil {
		message = fmt.Sprintf("%s: %v", message, original)
	}

	return &HTTPError{
		Service:    service,
		Operation:  operation,
		StatusCode: statusCode,
		Message:    message,
		RequestURL: requestURL,
		Original:   original,
	}
}

// LogHTTPRequest 记录HTTP请求日志（仅在开发环境）
func LogHTTPRequest(service, operation, method, url string) {
	if os.Getenv("GO_ENV") == "development" {
		log.Printf("[%s:%s] %s %s", service, operation, method, url)
	}
}

// LogHTTPResponse 记录HTTP响应日志（仅在开发环境）
func LogHTTPResponse(service, operation string, resp *resty.Response, err error) {
	if os.Getenv("GO_ENV") != "development" {
		return
	}

	if err != nil {
		log.Printf("[%s:%s] Request failed: %v", service, operation, err)
		return
	}

	if resp != nil {
		log.Printf("[%s:%s] Response: %d %s (Duration: %v)",
			service, operation, resp.StatusCode(), resp.Status(), resp.Time())

		if resp.IsError() {
			log.Printf("[%s:%s] Error body: %s", service, operation, string(resp.Body()))
		}
	}
}

// HandleHTTPError 统一处理HTTP错误
func HandleHTTPError(service, operation string, resp *resty.Response, err error) error {
	// 记录响应日志
	LogHTTPResponse(service, operation, resp, err)

	if err != nil {
		return NewHTTPError(service, operation, resp, err)
	}

	if resp != nil && resp.IsError() {
		return NewHTTPError(service, operation, resp, nil)
	}

	return nil
}