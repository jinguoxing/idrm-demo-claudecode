// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"fmt"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// 数据库配置
	DB struct {
		Host         string `json:",default=localhost"`
		Port         int    `json:",default=3306"`
		Database     string `json:",default=idrm"`
		Username     string `json:",default=root"`
		Password     string `json:",default=123456"`
		Charset      string `json:",default=utf8mb4"`
		MaxIdleConns int    `json:",default=10"`
		MaxOpenConns int    `json:",default=100"`
	} `json:",optional"`
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.DB.Username,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Database,
		c.DB.Charset,
	)
}
