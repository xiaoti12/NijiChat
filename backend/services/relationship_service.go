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

// RelationshipService 声优关系服务
type RelationshipService struct {
	db            *database.D1Client
	cache         *database.KVClient
	seiyuuService *SeiyuuService
	aiService     *AIService
}

// NewRelationshipService 创建关系服务实例
func NewRelationshipService(db *database.D1Client, cache *database.KVClient, seiyuuService *SeiyuuService, aiService *AIService) *RelationshipService {
	return &RelationshipService{
		db:            db,
		cache:         cache,
		seiyuuService: seiyuuService,
		aiService:     aiService,
	}
}

// GenerateRelationship 基于两个声优的原始资料AI生成关系
func (s *RelationshipService) GenerateRelationship(ctx context.Context, seiyuuIdA, seiyuuIdB string) (*models.GeneratedRelationship, error) {
	// 获取两个声优的信息
	seiyuuA, err := s.seiyuuService.GetSeiyuuByID(ctx, seiyuuIdA)
	if err != nil {
		return nil, fmt.Errorf("failed to get seiyuu A: %w", err)
	}

	seiyuuB, err := s.seiyuuService.GetSeiyuuByID(ctx, seiyuuIdB)
	if err != nil {
		return nil, fmt.Errorf("failed to get seiyuu B: %w", err)
	}

	// 调用AI服务生成关系
	generatedRel, err := s.aiService.GenerateSeiyuuRelationship(ctx, seiyuuA, seiyuuB)
	if err != nil {
		return nil, fmt.Errorf("failed to generate relationship with AI: %w", err)
	}

	return generatedRel, nil
}

// CreateRelationship 创建新关系
func (s *RelationshipService) CreateRelationship(ctx context.Context, req *models.CreateRelationshipRequest) (*models.SeiyuuRelationship, error) {
	// 生成UUID
	id := uuid.New().String()

	// 创建关系对象
	relationship := &models.SeiyuuRelationship{
		ID:                      id,
		SeiyuuIdA:               req.SeiyuuIdA,
		SeiyuuIdB:               req.SeiyuuIdB,
		RelationshipDescription: req.RelationshipDescription,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	// 标准化声优顺序（避免重复关系）
	relationship.NormalizeSeiyuuOrder()

	// 验证数据
	if err := relationship.Validate(); err != nil {
		return nil, err
	}

	// 检查声优是否存在
	if _, err := s.seiyuuService.GetSeiyuuByID(ctx, relationship.SeiyuuIdA); err != nil {
		return nil, fmt.Errorf("seiyuu A not found: %w", err)
	}
	if _, err := s.seiyuuService.GetSeiyuuByID(ctx, relationship.SeiyuuIdB); err != nil {
		return nil, fmt.Errorf("seiyuu B not found: %w", err)
	}

	// 插入数据库
	query := `
		INSERT INTO seiyuu_relationships (id, seiyuu_id_a, seiyuu_id_b, relationship_description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(ctx, query,
		relationship.ID,
		relationship.SeiyuuIdA,
		relationship.SeiyuuIdB,
		relationship.RelationshipDescription,
		database.TimeToString(relationship.CreatedAt),
		database.TimeToString(relationship.UpdatedAt),
	)

	if err != nil {
		// 检查是否是唯一约束冲突
		if database.IsUniqueConstraintError(err) {
			return nil, models.ErrRelationshipAlreadyExists
		}
		return nil, fmt.Errorf("failed to create relationship: %w", err)
	}

	return relationship, nil
}

// GetRelationship 获取两个声优之间的关系
func (s *RelationshipService) GetRelationship(ctx context.Context, seiyuuIdA, seiyuuIdB string) (*models.SeiyuuRelationship, error) {
	// 标准化顺序
	if seiyuuIdA > seiyuuIdB {
		seiyuuIdA, seiyuuIdB = seiyuuIdB, seiyuuIdA
	}

	query := `
		SELECT id, seiyuu_id_a, seiyuu_id_b, relationship_description, created_at, updated_at
		FROM seiyuu_relationships
		WHERE seiyuu_id_a = ? AND seiyuu_id_b = ?
	`

	row := s.db.QueryRow(ctx, query, seiyuuIdA, seiyuuIdB)
	relationship, err := s.scanRelationshipRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrRelationshipNotFound
		}
		return nil, fmt.Errorf("failed to query relationship: %w", err)
	}

	return relationship, nil
}

// GetRelationshipByID 根据ID获取关系
func (s *RelationshipService) GetRelationshipByID(ctx context.Context, id string) (*models.SeiyuuRelationship, error) {
	query := `
		SELECT id, seiyuu_id_a, seiyuu_id_b, relationship_description, created_at, updated_at
		FROM seiyuu_relationships
		WHERE id = ?
	`

	row := s.db.QueryRow(ctx, query, id)
	relationship, err := s.scanRelationshipRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrRelationshipNotFound
		}
		return nil, fmt.Errorf("failed to query relationship: %w", err)
	}

	return relationship, nil
}

// UpdateRelationship 更新关系信息
func (s *RelationshipService) UpdateRelationship(ctx context.Context, id string, req *models.UpdateRelationshipRequest) (*models.SeiyuuRelationship, error) {
	// 先获取现有关系
	relationship, err := s.GetRelationshipByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.RelationshipDescription != nil {
		relationship.RelationshipDescription = *req.RelationshipDescription
	}
	relationship.UpdatedAt = time.Now()

	// 验证数据
	if err := relationship.Validate(); err != nil {
		return nil, err
	}

	// 更新数据库
	query := `
		UPDATE seiyuu_relationships
		SET relationship_description = ?, updated_at = ?
		WHERE id = ?
	`

	_, err = s.db.Exec(ctx, query,
		relationship.RelationshipDescription,
		database.TimeToString(relationship.UpdatedAt),
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update relationship: %w", err)
	}

	return relationship, nil
}

// DeleteRelationship 删除关系
func (s *RelationshipService) DeleteRelationship(ctx context.Context, id string) error {
	query := `DELETE FROM seiyuu_relationships WHERE id = ?`
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete relationship: %w", err)
	}

	// 检查是否有行被删除
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrRelationshipNotFound
	}

	return nil
}

// GetSeiyuuRelationships 获取某个声优的所有关系
func (s *RelationshipService) GetSeiyuuRelationships(ctx context.Context, seiyuuId string) ([]*models.SeiyuuRelationshipWithNames, error) {
	query := `
		SELECT
			r.id, r.seiyuu_id_a, r.seiyuu_id_b, r.relationship_description, r.created_at, r.updated_at,
			sa.name as seiyuu_name_a, sb.name as seiyuu_name_b
		FROM seiyuu_relationships r
		JOIN seiyuu sa ON r.seiyuu_id_a = sa.id
		JOIN seiyuu sb ON r.seiyuu_id_b = sb.id
		WHERE r.seiyuu_id_a = ? OR r.seiyuu_id_b = ?
		ORDER BY r.created_at DESC
	`

	rows, err := s.db.Query(ctx, query, seiyuuId, seiyuuId)
	if err != nil {
		return nil, fmt.Errorf("failed to query seiyuu relationships: %w", err)
	}
	defer rows.Close()

	return s.scanRelationshipWithNamesRows(rows)
}

// GetAllRelationships 获取所有关系（管理员接口）
func (s *RelationshipService) GetAllRelationships(ctx context.Context) ([]*models.SeiyuuRelationshipWithNames, error) {
	query := `
		SELECT
			r.id, r.seiyuu_id_a, r.seiyuu_id_b, r.relationship_description, r.created_at, r.updated_at,
			sa.name as seiyuu_name_a, sb.name as seiyuu_name_b
		FROM seiyuu_relationships r
		JOIN seiyuu sa ON r.seiyuu_id_a = sa.id
		JOIN seiyuu sb ON r.seiyuu_id_b = sb.id
		ORDER BY r.created_at DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all relationships: %w", err)
	}
	defer rows.Close()

	return s.scanRelationshipWithNamesRows(rows)
}

// scanRelationshipRow 扫描单行关系数据
func (s *RelationshipService) scanRelationshipRow(row *sql.Row) (*models.SeiyuuRelationship, error) {
	var relationship models.SeiyuuRelationship
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&relationship.ID,
		&relationship.SeiyuuIdA,
		&relationship.SeiyuuIdB,
		&relationship.RelationshipDescription,
		&createdAtStr,
		&updatedAtStr,
	)

	if err != nil {
		return nil, err
	}

	// 转换时间字符串为time.Time
	relationship.CreatedAt, err = database.StringToTime(createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	relationship.UpdatedAt, err = database.StringToTime(updatedAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	return &relationship, nil
}

// scanRelationshipWithNamesRows 扫描多行关系数据（带声优姓名）
func (s *RelationshipService) scanRelationshipWithNamesRows(rows *sql.Rows) ([]*models.SeiyuuRelationshipWithNames, error) {
	relationships := make([]*models.SeiyuuRelationshipWithNames, 0)

	for rows.Next() {
		var relationship models.SeiyuuRelationshipWithNames
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&relationship.ID,
			&relationship.SeiyuuIdA,
			&relationship.SeiyuuIdB,
			&relationship.RelationshipDescription,
			&createdAtStr,
			&updatedAtStr,
			&relationship.SeiyuuNameA,
			&relationship.SeiyuuNameB,
		)

		if err != nil {
			return nil, err
		}

		// 转换时间字符串为time.Time
		relationship.CreatedAt, err = database.StringToTime(createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		relationship.UpdatedAt, err = database.StringToTime(updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}

		relationships = append(relationships, &relationship)
	}

	return relationships, nil
}