package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"seiyuu-chat/models"
)

// AIService AI服务（调用外部AI API）
type AIService struct {
	apiKey     string
	apiBaseURL string
	httpClient *http.Client
}

// NewAIService 创建AI服务实例
func NewAIService() *AIService {
	apiKey := os.Getenv("LIGHTWEIGHT_AI_API_KEY")
	if apiKey == "" {
		apiKey = "dev-ai-key" // 开发环境默认值
	}

	return &AIService{
		apiKey:     apiKey,
		apiBaseURL: "https://api.openai.com/v1", // 示例，实际应该从环境变量读取
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CallLightweightAI 调用轻量级AI模型（用于调度和资料处理）
func (s *AIService) CallLightweightAI(ctx context.Context, prompt string) (string, error) {
	// 构建请求
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

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiBaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API error: %s", string(body))
	}

	// 解析响应
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
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
