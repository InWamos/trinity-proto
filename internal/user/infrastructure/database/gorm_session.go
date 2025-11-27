package database

import "gorm.io/gorm"

type GormSession struct {
	tx *gorm.DB
}
