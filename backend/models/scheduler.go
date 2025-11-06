package models

// SchedulerRequest 智能调度器请求
type SchedulerRequest struct {
	Message          string   `json:"message" binding:"required"`      // 用户消息
	GroupID          string   `json:"group_id,omitempty"`              // 群组ID（可选）
	AvailableSeiyuu  []string `json:"available_seiyuu,omitempty"`      // 可选声优列表（可选）
	Context          string   `json:"context,omitempty"`               // 上下文对话历史
}

// SchedulerResponse 智能调度器响应
type SchedulerResponse struct {
	Success bool              `json:"success"`
	Data    *SchedulerResult  `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}

// SchedulerResult 调度结果
type SchedulerResult struct {
	SelectedSeiyuu *Seiyuu `json:"selected_seiyuu"` // 选中的声优
	Reason         string  `json:"reason"`          // 选择原因
	Confidence     float64 `json:"confidence"`      // 置信度 (0-1)
}

// MoegirlRawDataRequest 萌娘百科原始数据请求
type MoegirlRawDataRequest struct {
	Name string `json:"name" binding:"required"` // 声优名称
}

// MoegirlRawDataResponse 萌娘百科原始数据响应
type MoegirlRawDataResponse struct {
	Success bool                `json:"success"`
	Data    *MoegirlRawData     `json:"data,omitempty"`
	Error   string              `json:"error,omitempty"`
}

// MoegirlRawData 萌娘百科原始数据
type MoegirlRawData struct {
	RawText   string `json:"raw_text"`    // 原始wiki文本
	PageTitle string `json:"page_title"`  // 页面标题
	PageURL   string `json:"page_url"`    // 页面URL
}

// ProcessProfileRequest AI处理资料请求
type ProcessProfileRequest struct {
	RawText     string `json:"raw_text" binding:"required"`      // 原始文本
	SeiyuuName  string `json:"seiyuu_name" binding:"required"`   // 声优名称
}

// ProcessProfileResponse AI处理资料响应
type ProcessProfileResponse struct {
	Success bool                `json:"success"`
	Data    *ProcessedProfile   `json:"data,omitempty"`
	Error   string              `json:"error,omitempty"`
}

// ProcessedProfile 处理后的资料
type ProcessedProfile struct {
	ProfileMarkdown string   `json:"profile_markdown"` // 处理后的Markdown格式资料
	SuggestedTags   []string `json:"suggested_tags"`   // 建议的标签
}
