package services

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"seiyuu-chat/models"
	"seiyuu-chat/utils"

	"github.com/go-resty/resty/v2"
)

// MoegirlService 萌娘百科服务
type MoegirlService struct {
	apiBaseURL  string
	restyClient *resty.Client
}

// NewMoegirlService 创建萌娘百科服务实例
func NewMoegirlService() *MoegirlService {
	apiBaseURL := os.Getenv("MOEGIRL_API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://zh.moegirl.org.cn/api.php"
	}

	// 使用共享的HTTP客户端配置
	config := utils.DefaultHTTPConfig()
	// 为萌娘百科API定制超时时间
	config.Timeout = 15 * time.Second
	config.RetryMaxWaitTime = 3 * time.Second

	client := utils.NewRestyClient(config)

	return &MoegirlService{
		apiBaseURL:  apiBaseURL,
		restyClient: client,
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

	// 记录请求日志
	utils.LogHTTPRequest("moegirl", "GetRawData", "GET", s.apiBaseURL)

	// 发送HTTP请求
	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"action":  "query",
			"format":  "json",
			"prop":    "revisions",
			"titles":  name,
			"rvprop":  "content",
			"rvslots": "main",
		}).
		SetResult(&apiResponse).
		Get(s.apiBaseURL)

	// 统一错误处理
	if httpErr := utils.HandleHTTPError("moegirl", "GetRawData", resp, err); httpErr != nil {
		return nil, httpErr
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

	// 记录请求日志
	utils.LogHTTPRequest("moegirl", "SearchSeiyuu", "GET", s.apiBaseURL)

	// 发送HTTP请求
	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"action": "opensearch",
			"format": "json",
			"search": keyword,
			"limit":  "10",
		}).
		SetResult(&searchResults).
		Get(s.apiBaseURL)

	// 统一错误处理
	if httpErr := utils.HandleHTTPError("moegirl", "SearchSeiyuu", resp, err); httpErr != nil {
		return nil, httpErr
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
