package user

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// gormDao GORM 实现
type gormDao struct {
	db *gorm.DB
}

// NewGormDao 创建 GORM DAO
func NewGormDao(db *gorm.DB) Model {
	return &gormDao{db: db}
}

// Insert 插入用户
func (d *gormDao) Insert(ctx context.Context, data *User) (*User, error) {
	if err := d.db.WithContext(ctx).Create(data).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}
	return data, nil
}

// FindByPhone 根据手机号查找用户
func (d *gormDao) FindByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	err := d.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (d *gormDao) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := d.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// FindOne 根据 ID 查找用户
func (d *gormDao) FindOne(ctx context.Context, id string) (*User, error) {
	var user User
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// Update 更新用户信息
func (d *gormDao) Update(ctx context.Context, data *User) error {
	if err := d.db.WithContext(ctx).Save(data).Error; err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	return nil
}

// UpdateLastLogin 更新最后登录时间
func (d *gormDao) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	if err := d.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("last_login_at", now).Error; err != nil {
		return fmt.Errorf("更新最后登录时间失败: %w", err)
	}
	return nil
}

// WithTx 设置事务
func (d *gormDao) WithTx(tx interface{}) Model {
	if txDB, ok := tx.(*gorm.DB); ok {
		return &gormDao{db: txDB}
	}
	return d
}

// Trans 执行事务
func (d *gormDao) Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, &gormDao{db: tx})
	})
}
