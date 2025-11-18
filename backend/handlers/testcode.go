package handlers

import (
	"errors"
	"net/http"
	"seiyuu-chat/router"
	"seiyuu-chat/utils"
)

// TestCodeHandler 测试码功能API处理器
type TestCodeHandler struct{}

// NewTestCodeHandler 创建测试码处理器实例
func NewTestCodeHandler() *TestCodeHandler {
	return &TestCodeHandler{}
}

// ModelConfigResponse AI模型配置响应
type ModelConfigResponse struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	APIKey        string  `json:"api_key"`
	APIEndpoint   string  `json:"api_endpoint"`
	ModelName     string  `json:"model_name"`
	IsLightweight bool    `json:"is_lightweight"`
	MaxTokens     int     `json:"max_tokens"`
	Temperature   float64 `json:"temperature"`
}

// 测试码映射配置（硬编码）
var testCodeConfigs = map[string]ModelConfigResponse{
	"nijichat": {
		Name:          "体验模型",
		Type:          "openai",
		APIEndpoint:   "https://yunwu.ai",
		ModelName:     "deepseek-chat",
		IsLightweight: false,
		MaxTokens:     2000,
		Temperature:   0.7,
	},
}

// GetModelConfig 通过测试码获取AI模型配置
// GET /api/test-model-config?code=xxx
func (h *TestCodeHandler) GetModelConfig(c *router.Context) {
	// 获取查询参数中的测试码
	code := c.Query("code")
	if code == "" {
		utils.BadRequestError(c, errors.New("缺少测试码参数"))
		return
	}

	// 查找测试码对应的配置
	config, exists := testCodeConfigs[code]
	if !exists {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "无效的测试码",
		})
		return
	}

	// 从环境变量获取API Key
	apiKey := utils.GetEnv("TEST_MODEL_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "服务器未配置测试API Key",
		})
		return
	}

	// 填充API Key
	config.APIKey = apiKey

	// 返回配置
	utils.SuccessResponse(c, config)
}
