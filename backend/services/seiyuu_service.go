package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"seiyuu-chat/database"
	"seiyuu-chat/models"

	"github.com/google/uuid"
)

// SeiyuuService 声优服务
type SeiyuuService struct {
	db    *database.D1Client
	cache *database.KVClient
}

// NewSeiyuuService 创建声优服务实例
func NewSeiyuuService(db *database.D1Client, cache *database.KVClient) *SeiyuuService {
	return &SeiyuuService{
		db:    db,
		cache: cache,
	}
}

// GetAllSeiyuu 获取所有已发布的声优列表
func (s *SeiyuuService) GetAllSeiyuu(ctx context.Context) ([]*models.Seiyuu, error) {
	query := `
		SELECT id, name, avatar_url, profile_markdown, tags, status, created_at, updated_at
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
		SELECT id, name, avatar_url, profile_markdown, tags, status, created_at, updated_at
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

	// 创建声优对象
	seiyuu := &models.Seiyuu{
		ID:              id,
		Name:            req.Name,
		AvatarURL:       req.AvatarURL,
		ProfileMarkdown: req.ProfileMarkdown,
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
		INSERT INTO seiyuu (id, name, avatar_url, profile_markdown, tags, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	tagsJSON := database.JSONArrayToString(seiyuu.Tags)
	_, err := s.db.Exec(ctx, query,
		seiyuu.ID,
		seiyuu.Name,
		seiyuu.AvatarURL,
		seiyuu.ProfileMarkdown,
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
		SET name = ?, avatar_url = ?, profile_markdown = ?, tags = ?, status = ?, updated_at = ?
		WHERE id = ?
	`

	tagsJSON := database.JSONArrayToString(seiyuu.Tags)
	_, err = s.db.Exec(ctx, query,
		seiyuu.Name,
		seiyuu.AvatarURL,
		seiyuu.ProfileMarkdown,
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
		SELECT id, name, avatar_url, profile_markdown, tags, status, created_at, updated_at
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
	var avatarURL sql.NullString
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&seiyuu.ID,
		&seiyuu.Name,
		&avatarURL,
		&seiyuu.ProfileMarkdown,
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
		var avatarURL sql.NullString
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&seiyuu.ID,
			&seiyuu.Name,
			&avatarURL,
			&seiyuu.ProfileMarkdown,
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
