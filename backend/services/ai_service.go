package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"seiyuu-chat/models"
	"seiyuu-chat/utils"

	"github.com/syumai/workers/cloudflare"
)

// AIService AI服务（调用外部AI API）
type AIService struct {
	apiKey     string
	apiBaseURL string
	modelName  string
	httpClient *http.Client
}

// NewAIService 创建AI服务实例
func NewAIService() (*AIService, error) {
	apiKey := cloudflare.Getenv("LIGHTWEIGHT_AI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("LIGHTWEIGHT_AI_API_KEY 环境变量未设置")
	}

	apiBaseURL := cloudflare.Getenv("AI_API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://api.openai.com/v1" // 默认使用OpenAI API
	}

	modelName := cloudflare.Getenv("AI_MODEL_NAME")
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}

	// 使用 Cloudflare Workers 兼容的 HTTP 客户端
	config := utils.DefaultHTTPConfig()
	cloudflareClient := utils.NewCloudflareHTTPClientWithAuth(config, apiKey)

	return &AIService{
		apiKey:     apiKey,
		apiBaseURL: apiBaseURL,
		modelName:  modelName,
		httpClient: cloudflareClient.HTTPClient(),
	}, nil
}

// CallLightweightAI 调用轻量级AI模型（用于调度和资料处理）
func (s *AIService) CallLightweightAI(ctx context.Context, prompt string) (string, error) {
	// 构建请求体
	requestBody := map[string]interface{}{
		"model": s.modelName, // 使用配置的模型名称
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.2,
		"max_tokens":  2000,
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

	// 序列化请求体
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 记录请求日志
	apiURL := s.apiBaseURL + "/chat/completions"
	utils.LogHTTPRequest("ai", "CallLightweightAI", "POST", apiURL)

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NijiChat/1.0")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用AI API失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AI API响应错误，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 读取并解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析AI响应失败: %w", err)
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
1. 总结关键信息：姓名、生日、代表作品、个人爱好、人际关系、个人轶事等
2. 组织为清晰的Markdown格式
3. 保持客观真实，不添加虚构内容，不允许增加或删除信息，只能整理已有内容
4. 建议3-5个相关标签（例如性格、爱好等）

原始资料：
%s

请直接返回纯JSON格式，不要使用markdown代码块：
重要：不要添加代码块标记或任何其他格式，直接输出可解析的JSON对象！
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

// GenerateSeiyuuRelationship 使用AI生成声优关系
func (s *AIService) GenerateSeiyuuRelationship(ctx context.Context, seiyuuA, seiyuuB *models.Seiyuu) (*models.GeneratedRelationship, error) {
	prompt := fmt.Sprintf(`请根据以下两位声优的资料分析他们之间可能存在的关系。

要求：
1. 基于提供的真实资料进行分析，不要添加虚构内容
2. 从中立的第三者视角描述关系，使用客观语言
3. 详细说明两人的具体关系背景、共同点或互动情况
4. 如果没有明显关系，描述他们作为同行的共同特点和专业领域

声优A：%s
原始资料：
%s

声优B：%s
原始资料：
%s

请直接返回纯JSON格式，不要使用markdown代码块：
重要：不要添加代码块标记或任何其他格式，直接输出可解析的JSON对象！
{
  "relationship_description": "详细的关系描述，从上帝视角客观阐述两人的关系背景、共同点或互动情况"
}`, seiyuuA.Name, seiyuuA.RawProfileData, seiyuuB.Name, seiyuuB.RawProfileData)

	response, err := s.CallLightweightAI(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// 解析JSON响应
	var result models.GeneratedRelationship
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		// 如果解析失败，返回默认关系
		return &models.GeneratedRelationship{
			RelationshipDescription: fmt.Sprintf("%s 和 %s 同为声优行业的从业者，各自在不同的作品和角色中展现才华。", seiyuuA.Name, seiyuuB.Name),
		}, nil
	}

	// 验证生成的内容不为空
	if result.RelationshipDescription == "" {
		result.RelationshipDescription = fmt.Sprintf("%s 和 %s 同为声优行业的从业者。", seiyuuA.Name, seiyuuB.Name)
	}

	return &result, nil
}
