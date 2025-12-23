package royalties

import (
	"log"

	"github.com/arod1213/auto_ingestion/models"
	"gorm.io/gorm"
)

func (p ExtPayment) ToPayment(songID *uint, userID uint) Payment {
	return Payment{
		UserID:   userID,
		SongID:   songID,
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
			err := db.Save(&payments).Error
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
	var song models.Song

	if p.Isrc != nil {
		if v, ok := cache[*p.Isrc]; ok {
			payment := p.ToPayment(&v, userID)
			return &payment, nil
		}

		err := db.Where("isrc = ?", *p.Isrc).First(&song).Error
		if err != nil {
			return nil, err
		}

		cache[*p.Isrc] = song.ID
	} else if p.Iswc != nil {
		if v, ok := cache[*p.Iswc]; ok {
			payment := p.ToPayment(&v, userID)
			return &payment, nil
		}

		err := db.Where("iswc = ?", *p.Iswc).First(&song).Error
		if err != nil {
			return nil, err
		}

		cache[*p.Iswc] = song.ID
	} else {
		if v, ok := cache[p.Title]; ok {
			payment := p.ToPayment(&v, userID)
			return &payment, nil
		}

		query := db.
			Joins("LEFT JOIN shares on shares.song_id = songs.id").
			Where("songs.title LIKE ?", "%"+p.Title+"%")

		if p.Artist != nil {
			query = query.Where("songs.artist LIKE ?", "%"+*p.Artist+"%")
		}

		err := query.
			Where("shares.user_id = ?", userID).
			First(&song).
			Error

		if err != nil {
			return nil, err
		}
		cache[p.Title] = song.ID
	}

	payment := p.ToPayment(&song.ID, userID)
	if song.ID == 0 {
		payment.SongID = nil
	}

	return &payment, nil
}
