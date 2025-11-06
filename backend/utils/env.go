package utils

import (
	"os"
	"strings"
)

// IsDevelopmentMode 检查当前是否为开发环境
func IsDevelopmentMode() bool {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	return env == "development" || env == "dev" || env == ""
}

// IsProductionMode 检查当前是否为生产环境
func IsProductionMode() bool {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	return env == "production" || env == "prod"
}

// IsStagingMode 检查当前是否为测试环境
func IsStagingMode() bool {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	return env == "staging" || env == "test"
}

// GetEnvironment 获取当前环境名称
func GetEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		return "development"
	}
	return env
}