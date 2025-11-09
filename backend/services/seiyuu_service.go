package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"seiyuu-chat/database"
	"seiyuu-chat/models"
	"seiyuu-chat/utils"

	"github.com/google/uuid"
)

// SeiyuuService 声优服务
type SeiyuuService struct {
	db         *database.D1Client
	cache      *database.KVClient
	httpClient *http.Client
	moegirlURL string
}

// NewSeiyuuService 创建声优服务实例
func NewSeiyuuService(db *database.D1Client, cache *database.KVClient) *SeiyuuService {
	// 设置萌娘百科API URL
	moegirlURL := utils.GetEnv("MOEGIRL_API_BASE_URL")
	if moegirlURL == "" {
		moegirlURL = "https://zh.moegirl.org.cn/api.php"
	}

	// 使用 Cloudflare Workers 兼容的 HTTP 客户端
	config := utils.DefaultHTTPConfig()
	cloudflareClient := utils.NewCloudflareHTTPClient(config)

	return &SeiyuuService{
		db:         db,
		cache:      cache,
		httpClient: cloudflareClient.HTTPClient(),
		moegirlURL: moegirlURL,
	}
}

// GetAllSeiyuu 获取所有已发布的声优列表
func (s *SeiyuuService) GetAllSeiyuu(ctx context.Context) ([]*models.Seiyuu, error) {
	query := `
		SELECT id, name, avatar_url, profile_markdown, raw_profile_data, tags, status, created_at, updated_at
		FROM seiyuu
		WHERE status = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, query, models.SeiyuuStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to query seiyuu: %w", err)
	}
	defer rows.Close()

	return s.scanSeiyuuRows(rows)
}

// GetAllSeiyuuAdmin 获取所有声优列表（管理员接口，包含所有状态）
func (s *SeiyuuService) GetAllSeiyuuAdmin(ctx context.Context) ([]*models.Seiyuu, error) {
	query := `
		SELECT id, name, avatar_url, profile_markdown, raw_profile_data, tags, status, created_at, updated_at
		FROM seiyuu
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all seiyuu: %w", err)
	}
	defer rows.Close()

	return s.scanSeiyuuRows(rows)
}

// GetSeiyuuByID 根据ID获取声优详情
func (s *SeiyuuService) GetSeiyuuByID(ctx context.Context, id string) (*models.Seiyuu, error) {
	// 先尝试从缓存获取
	cacheKey := database.CacheKey(database.CachePrefixSeiyuu, id)
	var seiyuu models.Seiyuu
	err := s.cache.GetJSON(ctx, cacheKey, &seiyuu)
	if err == nil {
		return &seiyuu, nil
	}

	// 从数据库查询
	query := `
		SELECT id, name, avatar_url, profile_markdown, raw_profile_data, tags, status, created_at, updated_at
		FROM seiyuu
		WHERE id = ?
	`

	row := s.db.QueryRow(ctx, query, id)
	seiyuuPtr, err := s.scanSeiyuuRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrSeiyuuNotFound
		}
		return nil, fmt.Errorf("failed to query seiyuu: %w", err)
	}

	// 存入缓存
	_ = s.cache.PutJSON(ctx, cacheKey, seiyuuPtr, database.CacheTTLMedium)

	return seiyuuPtr, nil
}

// CreateSeiyuu 创建声优
func (s *SeiyuuService) CreateSeiyuu(ctx context.Context, req *models.CreateSeiyuuRequest) (*models.Seiyuu, error) {
	// 生成UUID
	id := uuid.New().String()

	// 如果未提供原始资料，则使用ProfileMarkdown作为默认值
	rawProfileData := req.RawProfileData
	if rawProfileData == "" {
		rawProfileData = req.ProfileMarkdown
	}

	// 创建声优对象
	seiyuu := &models.Seiyuu{
		ID:              id,
		Name:            req.Name,
		AvatarURL:       req.AvatarURL,
		ProfileMarkdown: req.ProfileMarkdown,
		RawProfileData:  rawProfileData,
		Tags:            req.Tags,
		Status:          models.SeiyuuStatusPending, // 默认待审核状态
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 验证数据
	if err := seiyuu.Validate(); err != nil {
		return nil, err
	}

	// 插入数据库
	query := `
		INSERT INTO seiyuu (id, name, avatar_url, profile_markdown, raw_profile_data, tags, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	tagsJSON := database.JSONArrayToString(seiyuu.Tags)
	_, err := s.db.Exec(ctx, query,
		seiyuu.ID,
		seiyuu.Name,
		seiyuu.AvatarURL,
		seiyuu.ProfileMarkdown,
		seiyuu.RawProfileData,
		tagsJSON,
		seiyuu.Status,
		database.TimeToString(seiyuu.CreatedAt),
		database.TimeToString(seiyuu.UpdatedAt),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create seiyuu: %w", err)
	}

	return seiyuu, nil
}

// UpdateSeiyuu 更新声优信息
func (s *SeiyuuService) UpdateSeiyuu(ctx context.Context, id string, req *models.UpdateSeiyuuRequest) (*models.Seiyuu, error) {
	// 先获取现有声优
	seiyuu, err := s.GetSeiyuuByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != nil {
		seiyuu.Name = *req.Name
	}
	if req.AvatarURL != nil {
		seiyuu.AvatarURL = *req.AvatarURL
	}
	if req.ProfileMarkdown != nil {
		seiyuu.ProfileMarkdown = *req.ProfileMarkdown
	}
	if req.RawProfileData != nil {
		seiyuu.RawProfileData = *req.RawProfileData
	}
	if req.Tags != nil {
		seiyuu.Tags = *req.Tags
	}
	if req.Status != nil {
		seiyuu.Status = *req.Status
	}
	seiyuu.UpdatedAt = time.Now()

	// 验证数据
	if err := seiyuu.Validate(); err != nil {
		return nil, err
	}

	// 更新数据库
	query := `
		UPDATE seiyuu
		SET name = ?, avatar_url = ?, profile_markdown = ?, raw_profile_data = ?, tags = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	tagsJSON := database.JSONArrayToString(seiyuu.Tags)
	_, err = s.db.Exec(ctx, query,
		seiyuu.Name,
		seiyuu.AvatarURL,
		seiyuu.ProfileMarkdown,
		seiyuu.RawProfileData,
		tagsJSON,
		seiyuu.Status,
		database.TimeToString(seiyuu.UpdatedAt),
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update seiyuu: %w", err)
	}

	// 清除缓存
	cacheKey := database.CacheKey(database.CachePrefixSeiyuu, id)
	_ = s.cache.Delete(ctx, cacheKey)

	return seiyuu, nil
}

// DeleteSeiyuu 删除声优
func (s *SeiyuuService) DeleteSeiyuu(ctx context.Context, id string) error {
	query := `DELETE FROM seiyuu WHERE id = ?`
	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete seiyuu: %w", err)
	}

	// 清除缓存
	cacheKey := database.CacheKey(database.CachePrefixSeiyuu, id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}

// GetSeiyuuByIDs 根据ID列表批量获取声优
func (s *SeiyuuService) GetSeiyuuByIDs(ctx context.Context, ids []string) ([]*models.Seiyuu, error) {
	if len(ids) == 0 {
		return []*models.Seiyuu{}, nil
	}

	// 构建IN查询
	query := `
		SELECT id, name, avatar_url, profile_markdown, raw_profile_data, tags, status, created_at, updated_at
		FROM seiyuu
		WHERE id IN (?` + string(make([]byte, len(ids)-1)) + `)
	`

	// 转换为interface{}切片
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query seiyuu by IDs: %w", err)
	}
	defer rows.Close()

	return s.scanSeiyuuRows(rows)
}

// scanSeiyuuRow 扫描单行声优数据
func (s *SeiyuuService) scanSeiyuuRow(row *sql.Row) (*models.Seiyuu, error) {
	var seiyuu models.Seiyuu
	var tagsJSON string
	var avatarURL, rawProfileData sql.NullString
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&seiyuu.ID,
		&seiyuu.Name,
		&avatarURL,
		&seiyuu.ProfileMarkdown,
		&rawProfileData,
		&tagsJSON,
		&seiyuu.Status,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		seiyuu.AvatarURL = avatarURL.String
	}

	if rawProfileData.Valid {
		seiyuu.RawProfileData = rawProfileData.String
	}

	tags, err := database.StringToJSONArray(tagsJSON)
	if err != nil {
		seiyuu.Tags = []string{}
	} else {
		seiyuu.Tags = tags
	}

	// 转换时间字符串为time.Time
	seiyuu.CreatedAt, err = database.StringToTime(createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	seiyuu.UpdatedAt, err = database.StringToTime(updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return &seiyuu, nil
}

// scanSeiyuuRows 扫描多行声优数据
func (s *SeiyuuService) scanSeiyuuRows(rows *sql.Rows) ([]*models.Seiyuu, error) {
	seiyuus := make([]*models.Seiyuu, 0)

	for rows.Next() {
		var seiyuu models.Seiyuu
		var tagsJSON string
		var avatarURL, rawProfileData sql.NullString
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&seiyuu.ID,
			&seiyuu.Name,
			&avatarURL,
			&seiyuu.ProfileMarkdown,
			&rawProfileData,
			&tagsJSON,
			&seiyuu.Status,
			&createdAtStr,
			&updatedAtStr,
		)

		if err != nil {
			return nil, err
		}

		if avatarURL.Valid {
			seiyuu.AvatarURL = avatarURL.String
		}

		if rawProfileData.Valid {
			seiyuu.RawProfileData = rawProfileData.String
		}

		tags, err := database.StringToJSONArray(tagsJSON)
		if err != nil {
			seiyuu.Tags = []string{}
		} else {
			seiyuu.Tags = tags
		}

		// 转换时间字符串为time.Time
		seiyuu.CreatedAt, err = database.StringToTime(createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		seiyuu.UpdatedAt, err = database.StringToTime(updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}

		seiyuus = append(seiyuus, &seiyuu)
	}

	return seiyuus, nil
}

// GetMoegirlRawData 获取萌娘百科原始数据
func (s *SeiyuuService) GetMoegirlRawData(ctx context.Context, name string) (*models.MoegirlRawData, error) {
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
	u, err := url.Parse(s.moegirlURL)
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
