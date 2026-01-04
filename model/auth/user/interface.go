package user

import (
	"context"
)

// Model 用户数据访问接口
type Model interface {
	// Insert 插入用户
	Insert(ctx context.Context, data *User) (*User, error)

	// FindByPhone 根据手机号查找用户
	FindByPhone(ctx context.Context, phone string) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindOne 根据 ID 查找用户
	FindOne(ctx context.Context, id string) (*User, error)

	// Update 更新用户信息
	Update(ctx context.Context, data *User) error

	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx context.Context, id string) error

	// WithTx 设置事务
	WithTx(tx interface{}) Model

	// Trans 执行事务
	Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}
