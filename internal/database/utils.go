package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDatabase(logger logger.Interface, dbPath string) (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	return
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
