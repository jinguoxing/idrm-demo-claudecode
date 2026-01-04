// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"fmt"

	"github.com/jinguoxing/idrm-demo-claudecode/api/internal/config"
	"github.com/jinguoxing/idrm-demo-claudecode/model/auth/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	DB        *gorm.DB
	UserModel user.Model
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(c.GetDSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("获取数据库连接失败: %v", err))
	}
	sqlDB.SetMaxIdleConns(c.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.DB.MaxOpenConns)

	// 初始化 UserModel
	userModel := user.NewGormDao(db)

	return &ServiceContext{
		Config:    c,
		DB:        db,
		UserModel: userModel,
	}
}
