package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/syumai/workers/cloudflare"
)

// IsDevelopmentMode 检查当前是否为开发环境
func IsDevelopmentMode() bool {
	// 优先使用 cloudflare.Getenv，回退到 os.Getenv（用于本地开发）
	env := getEnvironmentVar()
	if env == "" {
		env = os.Getenv("ENVIRONMENT")
	}

	result := env == "development" || env == "dev" || env == ""

	// 添加调试日志
	fmt.Printf("IsDevelopmentMode check: env='%s', result=%t\n", env, result)

	return result
}

// IsProductionMode 检查当前是否为生产环境
func IsProductionMode() bool {
	env := getEnvironmentVar()
	return strings.ToLower(env) == "production" || strings.ToLower(env) == "prod"
}

// IsStagingMode 检查当前是否为测试环境
func IsStagingMode() bool {
	env := getEnvironmentVar()
	return strings.ToLower(env) == "staging" || strings.ToLower(env) == "test"
}

// getEnvironmentVar 统一的环境变量获取方法
func getEnvironmentVar() string {
	// 优先使用 cloudflare.Getenv，回退到 os.Getenv（用于本地开发）
	env := cloudflare.Getenv("ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENVIRONMENT")
	}
	return env
}

// GetEnv 统一的环境变量获取方法（公开函数）
// 优先使用 cloudflare.Getenv，回退到 os.Getenv（用于本地开发）
func GetEnv(key string) string {
	env := cloudflare.Getenv(key)
	if env == "" {
		env = os.Getenv(key)
	}
	return env
}
