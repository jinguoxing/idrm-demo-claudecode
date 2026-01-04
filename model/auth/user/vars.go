package user

import (
	"errors"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("用户不存在")
	// ErrUserPhoneExists 手机号已存在
	ErrUserPhoneExists = errors.New("手机号已被注册")
	// ErrUserEmailExists 邮箱已存在
	ErrUserEmailExists = errors.New("邮箱已被注册")
	// ErrUserDisabled 用户已禁用
	ErrUserDisabled = errors.New("用户已禁用")
)

// 用户状态常量
const (
	UserStatusNormal   = 1 // 正常
	UserStatusDisabled = 0 // 禁用
)
