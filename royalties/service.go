package royalties

import (
	"log"

	"github.com/arod1213/auto_ingestion/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p ExtPayment) ToPayment(shareID *uint) Payment {
	return Payment{
		Hash:     p.Hash,
		ShareID:  shareID,
		Earnings: p.Earnings,
		Payor:    p.Payor.Name,
		// Date:      p.Date,
		Territory: p.Territory,
	}
}

func SavePayments(db *gorm.DB, userID uint, list []ExtPayment) error {
	cache := make(map[string]uint)
	var payments []Payment

	for _, p := range list {
		payment, err := p.FindPayment(db, userID, cache)
		if err != nil {
			continue
		}
		payments = append(payments, *payment)
		if len(payments) >= 1000 {
			err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&payments).Error
			if err != nil {
				log.Println("failed to save payments", err.Error())
				return err
			}
			payments = []Payment{} // reset
		}
	}

	err := db.Save(&payments).Error
	if err != nil {
		log.Println("failed to save payments", err.Error())
		return err
	}
	return nil
}

func (p ExtPayment) FindPayment(db *gorm.DB, userID uint, cache map[string]uint) (*Payment, error) {
	var share models.Share

	if p.Isrc != nil {
		if v, ok := cache[*p.Isrc]; ok {
			payment := p.ToPayment(&v)
			return &payment, nil
		}

		err := db.
			Joins("songs on songs.id = shares.song_id").
			Where("songs.isrc = ?", *p.Isrc).
			First(&share).
			Error

		if err != nil {
			return nil, err
		}

		cache[*p.Isrc] = share.ID
	} else if p.Iswc != nil {
		if v, ok := cache[*p.Iswc]; ok {
			payment := p.ToPayment(&v)
			return &payment, nil
		}

		err := db.
			Joins("songs on songs.id = shares.song_id").
			Where("songs.iswc = ?", *p.Iswc).
			First(&share).
			Error

		if err != nil {
			return nil, err
		}

		cache[*p.Iswc] = share.ID
	} else {
		if v, ok := cache[p.Title]; ok {
			payment := p.ToPayment(&v)
			return &payment, nil
		}

		query := db.
			Joins("LEFT JOIN songs on shares.song_id = songs.id").
			Where("songs.title LIKE ?", "%"+p.Title+"%")

		if p.Artist != nil {
			query = query.Where("songs.artist LIKE ?", "%"+*p.Artist+"%")
		}

		err := query.
			Where("shares.user_id = ?", userID).
			First(&share).
			Error

		if err != nil {
			return nil, err
		}
		cache[p.Title] = share.ID
	}

	payment := p.ToPayment(&share.ID)
	if share.ID == 0 {
		payment.ShareID = nil
	}

	return &payment, nil
}
