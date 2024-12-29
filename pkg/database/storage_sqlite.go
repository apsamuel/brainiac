package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var SqliteClient *gorm.DB

func NewSqliteClient(config SqliteConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.Filename), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
