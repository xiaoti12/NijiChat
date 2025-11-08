package models


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

