package database

import (
	"context"

	"gorm.io/gorm"
)

type GormTransaction struct {
	tx *gorm.DB
}

func NewGormTransaction(database GormDatabase, ctx context.Context) *GormTransaction {
	return &GormTransaction{database.engine.WithContext(ctx).Begin()}
}
