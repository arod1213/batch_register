package database

import (
	"github.com/arod1213/auto_ingestion/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Song{}, &models.Share{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
