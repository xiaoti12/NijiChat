package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"seiyuu-chat/models"
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

	return &MoegirlService{
		apiBaseURL: apiBaseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetRawData 获取萌娘百科原始数据
func (s *MoegirlService) GetRawData(ctx context.Context, name string) (*models.MoegirlRawData, error) {
	// 构建API请求参数
	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("prop", "revisions")
	params.Add("titles", name)
	params.Add("rvprop", "content")
	params.Add("rvslots", "main")

	// 构建完整URL
	fullURL := fmt.Sprintf("%s?%s", s.apiBaseURL, params.Encode())

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "NijiChat/1.0")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("moegirl API error: status %d", resp.StatusCode)
	}

	// 解析响应
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

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
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
	// 构建搜索API请求参数
	params := url.Values{}
	params.Add("action", "opensearch")
	params.Add("format", "json")
	params.Add("search", keyword)
	params.Add("limit", "10")

	fullURL := fmt.Sprintf("%s?%s", s.apiBaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "NijiChat/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析搜索结果
	var searchResults []interface{}
	err = json.Unmarshal(body, &searchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
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
