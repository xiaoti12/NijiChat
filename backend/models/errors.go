package models

import (
	"errors"
)

// 数据模型相关错误定义

var (
	// 声优相关错误
	ErrSeiyuuNotFound        = errors.New("声优不存在")
	ErrInvalidSeiyuuName     = errors.New("声优姓名不能为空")
	ErrInvalidSeiyuuProfile  = errors.New("声优资料不能为空")
	ErrInvalidSeiyuuStatus   = errors.New("无效的声优状态")
	ErrSeiyuuAlreadyExists   = errors.New("声优已存在")

	// 群组相关错误
	ErrGroupNotFound            = errors.New("群组不存在")
	ErrInvalidGroupName         = errors.New("群组名称不能为空")
	ErrInvalidGroupMembers      = errors.New("群组成员不能为空")
	ErrInvalidGroupMembersCount = errors.New("群讨论模式至少需要3个成员")
	ErrGroupAlreadyExists       = errors.New("群组已存在")

	// 管理员相关错误
	ErrAdminNotFound          = errors.New("管理员不存在")
	ErrInvalidAdminUsername   = errors.New("用户名不能为空")
	ErrInvalidAdminPassword   = errors.New("密码不能为空")
	ErrAdminAlreadyExists     = errors.New("管理员已存在")
	ErrInvalidCredentials     = errors.New("用户名或密码错误")
	ErrUnauthorized           = errors.New("未授权访问")
	ErrInvalidToken           = errors.New("无效的令牌")
	ErrTokenExpired           = errors.New("令牌已过期")

	// 通用错误
	ErrInternalServer         = errors.New("服务器内部错误")
	ErrBadRequest             = errors.New("请求参数错误")
	ErrNotFound               = errors.New("资源不存在")
	ErrDatabaseError          = errors.New("数据库操作失败")
	ErrCacheError             = errors.New("缓存操作失败")
)
