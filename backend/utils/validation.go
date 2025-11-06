package utils

import (
	"errors"
	"regexp"
	"strings"
)

// 数据验证工具函数

var (
	// 用户名正则：4-20位字母数字下划线
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{4,20}$`)

	// URL正则（简单验证）
	urlRegex = regexp.MustCompile(`^https?://[^\s]+$`)
)

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("用户名不能为空")
	}
	if !usernameRegex.MatchString(username) {
		return errors.New("用户名格式错误：需要4-20位字母、数字或下划线")
	}
	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}
	if len(password) < 6 {
		return errors.New("密码长度至少6位")
	}
	if len(password) > 100 {
		return errors.New("密码长度不能超过100位")
	}
	return nil
}

// ValidateURL 验证URL格式
func ValidateURL(url string) error {
	if url == "" {
		return nil // URL可以为空
	}
	if !urlRegex.MatchString(url) {
		return errors.New("URL格式错误")
	}
	return nil
}

// ValidateSeiyuuName 验证声优名称
func ValidateSeiyuuName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("声优名称不能为空")
	}
	if len(name) > 100 {
		return errors.New("声优名称长度不能超过100字符")
	}
	return nil
}

// ValidateMarkdown 验证Markdown内容
func ValidateMarkdown(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("Markdown内容不能为空")
	}
	if len(content) > 50000 {
		return errors.New("Markdown内容长度不能超过50000字符")
	}
	return nil
}

// ValidateTags 验证标签列表
func ValidateTags(tags []string) error {
	if len(tags) > 20 {
		return errors.New("标签数量不能超过20个")
	}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return errors.New("标签不能为空")
		}
		if len(tag) > 20 {
			return errors.New("单个标签长度不能超过20字符")
		}
	}
	return nil
}

// ValidateGroupName 验证群组名称
func ValidateGroupName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("群组名称不能为空")
	}
	if len(name) > 50 {
		return errors.New("群组名称长度不能超过50字符")
	}
	return nil
}

// ValidateMemberIDs 验证成员ID列表
func ValidateMemberIDs(memberIDs []string) error {
	if len(memberIDs) == 0 {
		return errors.New("成员列表不能为空")
	}
	if len(memberIDs) > 50 {
		return errors.New("成员数量不能超过50个")
	}

	// 检查是否有重复ID
	seen := make(map[string]bool)
	for _, id := range memberIDs {
		if seen[id] {
			return errors.New("成员列表中存在重复ID")
		}
		seen[id] = true
	}

	return nil
}

// SanitizeString 清理字符串（去除首尾空格）
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// SanitizeTags 清理标签列表
func SanitizeTags(tags []string) []string {
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}
