package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"seiyuu-chat/models"
	"seiyuu-chat/utils"
)

// MoegirlService 萌娘百科服务
type MoegirlService struct {
	apiBaseURL string
	httpClient *http.Client
}

// NewMoegirlService 创建萌娘百科服务实例
func NewMoegirlService() *MoegirlService {
	apiBaseURL := os.Getenv("MOEGIRL_API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://zh.moegirl.org.cn/api.php"
	}

	// 使用 Cloudflare Workers 兼容的 HTTP 客户端
	config := utils.DefaultHTTPConfig()
	cloudflareClient := utils.NewCloudflareHTTPClient(config)

	return &MoegirlService{
		apiBaseURL: apiBaseURL,
		httpClient: cloudflareClient.HTTPClient(),
	}
}

// GetRawData 获取萌娘百科原始数据
func (s *MoegirlService) GetRawData(ctx context.Context, name string) (*models.MoegirlRawData, error) {
	// 定义响应结构
	var apiResponse struct {
		Query struct {
			Pages map[string]struct {
				Title     string `json:"title"`
				Revisions []struct {
					Slots struct {
						Main struct {
							Content string `json:"*"`
						} `json:"main"`
					} `json:"slots"`
				} `json:"revisions"`
			} `json:"pages"`
		} `json:"query"`
	}

	// 构建请求URL
	u, err := url.Parse(s.apiBaseURL)
	if err != nil {
		return nil, fmt.Errorf("无效的API URL: %w", err)
	}

	query := u.Query()
	query.Set("action", "query")
	query.Set("format", "json")
	query.Set("prop", "revisions")
	query.Set("titles", name)
	query.Set("rvprop", "content")
	query.Set("rvslots", "main")
	u.RawQuery = query.Encode()

	// 记录请求日志
	utils.LogHTTPRequest("moegirl", "GetRawData", "GET", u.String())

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NijiChat/1.0")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求萌娘百科API失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API响应错误，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 读取并解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("解析API响应失败: %w", err)
	}

	// 提取内容
	for pageID, page := range apiResponse.Query.Pages {
		if pageID == "-1" {
			return nil, fmt.Errorf("页面不存在: %s", name)
		}

		if len(page.Revisions) == 0 {
			return nil, fmt.Errorf("页面内容为空: %s", name)
		}

		content := page.Revisions[0].Slots.Main.Content
		pageURL := fmt.Sprintf("https://zh.moegirl.org.cn/%s", url.PathEscape(page.Title))

		return &models.MoegirlRawData{
			RawText:   content,
			PageTitle: page.Title,
			PageURL:   pageURL,
		}, nil
	}

	return nil, fmt.Errorf("未找到页面内容")
}

// SearchSeiyuu 搜索声优页面
func (s *MoegirlService) SearchSeiyuu(ctx context.Context, keyword string) ([]string, error) {
	// 定义搜索响应结构
	var searchResults []interface{}

	// 构建请求URL
	u, err := url.Parse(s.apiBaseURL)
	if err != nil {
		return nil, fmt.Errorf("无效的API URL: %w", err)
	}

	query := u.Query()
	query.Set("action", "opensearch")
	query.Set("format", "json")
	query.Set("search", keyword)
	query.Set("limit", "10")
	u.RawQuery = query.Encode()

	// 记录请求日志
	utils.LogHTTPRequest("moegirl", "SearchSeiyuu", "GET", u.String())

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NijiChat/1.0")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("搜索萌娘百科失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("搜索API响应错误，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 读取并解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	err = json.Unmarshal(body, &searchResults)
	if err != nil {
		return nil, fmt.Errorf("解析搜索响应失败: %w", err)
	}

	if len(searchResults) < 2 {
		return []string{}, nil
	}

	// 提取标题列表
	titles, ok := searchResults[1].([]interface{})
	if !ok {
		return []string{}, nil
	}

	result := make([]string, 0, len(titles))
	for _, title := range titles {
		if titleStr, ok := title.(string); ok {
			result = append(result, titleStr)
		}
	}

	return result, nil
}
