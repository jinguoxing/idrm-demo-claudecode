package user

import (
	"gorm.io/gorm"
)

// ModelType 模型类型
type ModelType int

const (
	// ModelTypeGORM GORM 模型
	ModelTypeGORM ModelType = iota
	// ModelTypeSQLx SQLx 模型
	ModelTypeSQLx
)

// NewModel 创建模型
func NewModel(modelType ModelType, db interface{}) Model {
	switch modelType {
	case ModelTypeGORM:
		if gormDB, ok := db.(*gorm.DB); ok {
			return NewGormDao(gormDB)
		}
	case ModelTypeSQLx:
		return NewSqlxModel(db)
	}
	return nil
}
