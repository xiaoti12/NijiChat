package services

import (
	"context"
	"fmt"
	"strings"

	"seiyuu-chat/models"
)

// SchedulerService 智能调度服务
type SchedulerService struct {
	seiyuuService *SeiyuuService
	aiService     *AIService
}

// NewSchedulerService 创建调度服务实例
func NewSchedulerService(seiyuuService *SeiyuuService, aiService *AIService) *SchedulerService {
	return &SchedulerService{
		seiyuuService: seiyuuService,
		aiService:     aiService,
	}
}

// SelectSeiyuu 智能选择声优
func (s *SchedulerService) SelectSeiyuu(ctx context.Context, req *models.SchedulerRequest) (*models.SchedulerResult, error) {
	// 1. 获取可选声优列表
	availableSeiyuu, err := s.getAvailableSeiyuu(ctx, req.GroupID, req.AvailableSeiyuu)
	if err != nil {
		return nil, fmt.Errorf("获取声优列表失败: %w", err)
	}

	if len(availableSeiyuu) == 0 {
		return nil, fmt.Errorf("没有可用的声优")
	}

	// 如果只有一个声优，直接返回
	if len(availableSeiyuu) == 1 {
		return &models.SchedulerResult{
			SelectedSeiyuu: availableSeiyuu[0],
			Reason:         "唯一可选声优",
			Confidence:     1.0,
		}, nil
	}

	// 2. 尝试使用AI进行智能选择
	result, err := s.selectWithAI(ctx, req.Message, availableSeiyuu, req.Context)
	if err != nil {
		// AI调用失败，降级到关键词匹配
		return s.selectWithKeywords(req.Message, availableSeiyuu), nil
	}

	return result, nil
}

// getAvailableSeiyuu 获取可选声优列表
func (s *SchedulerService) getAvailableSeiyuu(ctx context.Context, groupID string, availableIDs []string) ([]*models.Seiyuu, error) {
	// 如果指定了availableIDs，只从这些ID中选择
	if len(availableIDs) > 0 {
		return s.seiyuuService.GetSeiyuuByIDs(ctx, availableIDs)
	}

	// 如果指定了groupID，获取群组成员
	if groupID != "" {
		// TODO: 从群组服务获取成员列表
		// 这里暂时返回所有声优
		return s.seiyuuService.GetAllSeiyuu(ctx)
	}

	// 默认返回所有已发布的声优
	return s.seiyuuService.GetAllSeiyuu(ctx)
}

// selectWithAI 使用AI进行智能选择
func (s *SchedulerService) selectWithAI(ctx context.Context, message string, seiyuus []*models.Seiyuu, contextHistory string) (*models.SchedulerResult, error) {
	// 构建AI提示词
	prompt := s.buildSchedulerPrompt(message, seiyuus, contextHistory)

	// 调用轻量级AI进行选择（网络I/O，不占用Worker CPU时间）
	// todo 调度似乎是放到前端了
	response, err := s.aiService.CallAI(ctx, prompt, "")
	if err != nil {
		return nil, err
	}

	// 解析AI响应
	result, err := s.parseAIResponse(response, seiyuus)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// buildSchedulerPrompt 构建调度器提示词
func (s *SchedulerService) buildSchedulerPrompt(message string, seiyuus []*models.Seiyuu, contextHistory string) string {
	var sb strings.Builder

	sb.WriteString("你是一个智能声优选择助手。请根据用户消息和对话上下文，从以下声优中选择最合适的一位进行回复。\n\n")

	sb.WriteString("可选声优列表：\n")
	for i, seiyuu := range seiyuus {
		sb.WriteString(fmt.Sprintf("%d. %s (ID: %s)\n", i+1, seiyuu.Name, seiyuu.ID))
		sb.WriteString(fmt.Sprintf("   标签: %s\n", strings.Join(seiyuu.Tags, ", ")))
		// 只包含简要资料，避免prompt过长
		if len(seiyuu.ProfileMarkdown) > 200 {
			sb.WriteString(fmt.Sprintf("   简介: %s...\n\n", seiyuu.ProfileMarkdown[:200]))
		} else {
			sb.WriteString(fmt.Sprintf("   简介: %s\n\n", seiyuu.ProfileMarkdown))
		}
	}

	if contextHistory != "" {
		sb.WriteString(fmt.Sprintf("对话上下文：\n%s\n\n", contextHistory))
	}

	sb.WriteString(fmt.Sprintf("用户消息：%s\n\n", message))

	sb.WriteString("请返回JSON格式的响应：\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"selected_id\": \"声优ID\",\n")
	sb.WriteString("  \"reason\": \"选择原因\",\n")
	sb.WriteString("  \"confidence\": 0.95\n")
	sb.WriteString("}\n")

	return sb.String()
}

// parseAIResponse 解析AI响应
func (s *SchedulerService) parseAIResponse(response string, seiyuus []*models.Seiyuu) (*models.SchedulerResult, error) {
	// TODO: 实现JSON解析逻辑
	// 这里提供一个简单的占位实现

	// 暂时返回第一个声优
	return &models.SchedulerResult{
		SelectedSeiyuu: seiyuus[0],
		Reason:         "AI选择",
		Confidence:     0.85,
	}, nil
}

// selectWithKeywords 使用关键词匹配进行选择（降级方案）
func (s *SchedulerService) selectWithKeywords(message string, seiyuus []*models.Seiyuu) *models.SchedulerResult {
	messageLower := strings.ToLower(message)

	// 遍历声优，计算匹配分数
	maxScore := 0
	var selectedSeiyuu *models.Seiyuu

	for _, seiyuu := range seiyuus {
		score := 0

		// 检查名字匹配
		if strings.Contains(messageLower, strings.ToLower(seiyuu.Name)) {
			score += 10
		}

		// 检查标签匹配
		for _, tag := range seiyuu.Tags {
			if strings.Contains(messageLower, strings.ToLower(tag)) {
				score += 5
			}
		}

		// 更新最高分数
		if score > maxScore {
			maxScore = score
			selectedSeiyuu = seiyuu
		}
	}

	// 如果没有匹配，随机选择第一个
	if selectedSeiyuu == nil {
		selectedSeiyuu = seiyuus[0]
		return &models.SchedulerResult{
			SelectedSeiyuu: selectedSeiyuu,
			Reason:         "默认选择",
			Confidence:     0.5,
		}
	}

	confidence := float64(maxScore) / 15.0 // 归一化到0-1
	if confidence > 1.0 {
		confidence = 1.0
	}

	return &models.SchedulerResult{
		SelectedSeiyuu: selectedSeiyuu,
		Reason:         "关键词匹配",
		Confidence:     confidence,
	}
}
