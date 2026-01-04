package user

import (
	"time"
	"gorm.io/gorm"
)

// User 用户实体
type User struct {
	Id           string         `gorm:"primaryKey;size:36"`             // UUID v7
	Phone        string         `gorm:"size:20;uniqueIndex:uk_phone"`   // 手机号
	Email        string         `gorm:"size:100;uniqueIndex:uk_email"`  // 邮箱
	PasswordHash string         `gorm:"size:255;not null"`              // 密码哈希
	Status       int8           `gorm:"type:tinyint;not null;default:1"` // 状态
	LastLoginAt  *time.Time     `gorm:"type:datetime"`                  // 最后登录时间
	CreatedAt    time.Time      `gorm:"type:datetime;not null"`
	UpdatedAt    time.Time      `gorm:"type:datetime;not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
