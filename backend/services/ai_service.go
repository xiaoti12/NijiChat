package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"seiyuu-chat/models"
	"seiyuu-chat/utils"

	"github.com/go-resty/resty/v2"
)

// AIService AI服务（调用外部AI API）
type AIService struct {
	apiKey      string
	apiBaseURL  string
	restyClient *resty.Client
}

// NewAIService 创建AI服务实例
func NewAIService() *AIService {
	apiKey := os.Getenv("LIGHTWEIGHT_AI_API_KEY")
	if apiKey == "" {
		apiKey = "dev-ai-key" // 开发环境默认值
	}

	apiBaseURL := os.Getenv("AI_API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://api.openai.com/v1" // 默认使用OpenAI API
	}

	// 使用共享的HTTP客户端配置
	config := utils.DefaultHTTPConfig()
	// 为AI API定制超时时间和重试配置
	config.Timeout = 30 * time.Second
	config.RetryMaxWaitTime = 5 * time.Second

	client := utils.NewRestyClientWithAuth(config, apiKey)

	return &AIService{
		apiKey:      apiKey,
		apiBaseURL:  apiBaseURL,
		restyClient: client,
	}
}

// CallLightweightAI 调用轻量级AI模型（用于调度和资料处理）
func (s *AIService) CallLightweightAI(ctx context.Context, prompt string) (string, error) {
	// 构建请求体
	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo", // 轻量级模型
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  500,
	}

	// 定义响应结构
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	// 记录请求日志
	apiURL := s.apiBaseURL + "/chat/completions"
	utils.LogHTTPRequest("ai", "CallLightweightAI", "POST", apiURL)

	// 发送HTTP请求
	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetBody(requestBody).
		SetResult(&response).
		Post(apiURL)

	// 统一错误处理
	if httpErr := utils.HandleHTTPError("ai", "CallLightweightAI", resp, err); httpErr != nil {
		return "", httpErr
	}

	// 检查API返回的业务错误
	if response.Error.Message != "" {
		return "", fmt.Errorf("AI API business error: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return response.Choices[0].Message.Content, nil
}

// ProcessSeiyuuProfile 使用AI处理声优资料
func (s *AIService) ProcessSeiyuuProfile(ctx context.Context, rawText string, seiyuuName string) (*models.ProcessedProfile, error) {
	prompt := fmt.Sprintf(`请将以下关于声优「%s」的原始资料整理为结构化的Markdown格式。

要求：
1. 提取关键信息：姓名、生日、血型、代表作品等
2. 组织为清晰的Markdown格式
3. 保持客观真实，不添加虚构内容
4. 建议3-5个相关标签（如"萝莉音"、"治愈系"等）

原始资料：
%s

请返回JSON格式：
{
  "profile_markdown": "整理后的Markdown文本",
  "suggested_tags": ["标签1", "标签2", "标签3"]
}`, seiyuuName, rawText)

	response, err := s.CallLightweightAI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// 解析JSON响应
	var result models.ProcessedProfile
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		// 如果解析失败，返回原始文本
		return &models.ProcessedProfile{
			ProfileMarkdown: response,
			SuggestedTags:   []string{},
		}, nil
	}

	return &result, nil
}
