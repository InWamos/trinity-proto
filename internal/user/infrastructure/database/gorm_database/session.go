package database

import "gorm.io/gorm"

type Session interface{}

type GormSession struct {
	tx *gorm.DB
}
