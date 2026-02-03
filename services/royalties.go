package services

import (
	"github.com/arod1213/auto_ingestion/models"
	"gorm.io/gorm"
)

type songEarnings struct {
	Song  models.Song `gorm:"embedded"`
	Total float32     `gorm:"column:total"`
}

type RoyaltyOverview struct {
	Total float32
	Songs []songEarnings
}

func GetRoyaltyOverview(db *gorm.DB, statementID uint) (*RoyaltyOverview, error) {
	var overview RoyaltyOverview
	var songs []songEarnings

	err := db.
		Table("payments p").
		Select("s.*, COALESCE(SUM(p.earnings), 0) AS total").
		Joins("JOIN songs s ON s.id = p.song_id").
		Where("p.statement_id = ?", statementID).
		Group("s.id").
		Scan(&songs).Error
	if err != nil {
		return nil, err
	}

	// Overall total
	err = db.
		Table("payments").
		Select("COALESCE(SUM(earnings), 0)").
		Where("statement_id = ?", statementID).
		Scan(&overview.Total).Error
	if err != nil {
		return nil, err
	}

	overview.Songs = songs
	return &overview, nil
}
