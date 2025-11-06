package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse 返回成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMessage 返回带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse 返回错误响应
func ErrorResponse(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   err.Error(),
	})
}

// ErrorWithMessage 返回带自定义消息的错误响应
func ErrorWithMessage(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// BadRequestError 返回400错误
func BadRequestError(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusBadRequest, err)
}

// UnauthorizedError 返回401错误
func UnauthorizedError(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusUnauthorized, err)
}

// ForbiddenError 返回403错误
func ForbiddenError(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusForbidden, err)
}

// NotFoundError 返回404错误
func NotFoundError(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusNotFound, err)
}

// InternalServerError 返回500错误
func InternalServerError(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusInternalServerError, err)
}

// PaginatedResponse 分页响应结构
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginatedResponse 创建分页响应
func NewPaginatedResponse(items interface{}, total, page, pageSize int) *PaginatedResponse {
	totalPages := (total + pageSize - 1) / pageSize
	return &PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// SuccessWithPagination 返回分页成功响应
func SuccessWithPagination(c *gin.Context, items interface{}, total, page, pageSize int) {
	response := NewPaginatedResponse(items, total, page, pageSize)
	SuccessResponse(c, response)
}
