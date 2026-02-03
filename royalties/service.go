package royalties

import (
	"errors"
	"log"

	"github.com/arod1213/auto_ingestion/models"
	"github.com/arod1213/auto_ingestion/royalties"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Reconcile(db *gorm.DB, userID uint) error {
	var payments []Payment
	err := db.Where("song_id is NULL AND user_id = ?", userID).Find(&payments).Error
	if err != nil {
		return err
	}

	cache := make(map[string]uint)
	for _, pay := range payments {
		payment, err := pay.Data.FindPayment(db, userID, cache)
		if err != nil {
			continue
		}

		pay.SongID = payment.SongID
		err = db.Save(&pay).Error
		if err != nil {
			log.Println("error saving payment")
			continue
		}
	}
	return nil
}

func (p ExtPayment) ToPayment(songID *uint, userID uint) Payment {
	return Payment{
		Data:     p,
		UserID:   userID,
		SongID:   songID,
		Hash:     p.Hash,
		Earnings: p.Earnings,
		Payor:    p.Payor.Name,
		// Date:      p.Date,
		Territory: p.Territory,
	}
}

func SavePayments(db *gorm.DB, userID uint, list []ExtPayment) (uint, error) {
	var s Statement
	err := db.Create(&s).Error
	if err != nil {
		log.Println("failed to create payment")
		return 0, err
	}

	cache := make(map[string]uint)
	var payments []Payment

	for _, p := range list {
		payment, err := p.FindPayment(db, userID, cache)
		if err != nil {
			continue
		}
		payment.StatementID = s.ID

		payments = append(payments, *payment)
		if len(payments) >= 1000 {
			err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&payments).Error
			if err != nil {
				log.Println("failed to save payments", err.Error())
				return 0, err
			}
			payments = []Payment{} // reset
		}
	}

	err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&payments).Error
	if err != nil {
		log.Println("failed to save payments", err.Error())
		return 0, err
	}
	return s.ID, nil
}

func (p ExtPayment) FindPayment(db *gorm.DB, userID uint, cache map[string]uint) (*Payment, error) {
	var song models.Song
	var songID *uint

	if p.Isrc != nil {
		if v, ok := cache[*p.Isrc]; ok {
			payment := p.ToPayment(&v, userID)
			return &payment, nil
		}

		err := db.
			Where("isrc = ?", *p.Isrc).
			First(&song).
			Error

		if err == nil {
			cache[*p.Isrc] = song.ID
			songID = &song.ID
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		cache[*p.Isrc] = song.ID
	} else if p.Iswc != nil {
		if v, ok := cache[*p.Iswc]; ok {
			payment := p.ToPayment(&v, userID)
			return &payment, nil
		}

		err := db.
			Where("iswc = ?", *p.Iswc).
			First(&song).
			Error

		if err == nil {
			cache[*p.Iswc] = song.ID
			songID = &song.ID
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
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

		if err == nil {
			cache[p.Title] = song.ID
			songID = &song.ID
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	payment := p.ToPayment(songID, userID)
	return &payment, nil
}
