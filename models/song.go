package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Society string

const (
	ASCAP  Society = "ASCAP"
	AMRA   Society = "AMRA"
	BMI    Society = "BMI"
	TheMLC Society = "The MLC"
	PRS    Society = "PRS"
	// HFA    Society = "HFA"
)

func (s Society) ToCode() string {
	switch s {
	case ASCAP:
		return "010"
	case AMRA:
		return "0"
	case BMI:
		return "021"
	case TheMLC:
		return "046"
	case PRS:
		return "052"
		// 	return "034"
		// case HFA:
	}
	panic("unreachable - fell through switch")
}

type Song struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Title  string `gorm:"type:varchar(255);not null"`
	Artist string `gorm:"type:varchar(255);not null"`
	Label  string `gorm:"type:varchar(255)"`

	Iswc *string `gorm:"type:varchar(15)"`
	Isrc string  `gorm:"type:varchar(15);uniqueIndex;not null"`
	Upc  uint64  `gorm:"type:integer;not null"`

	SpotifyID string `gorm:"type:text;not null"`
	// Url       string  `gorm:"type:text;not null"`
	Image *string `gorm:"type:text"`

	ReleaseDate time.Time     `gorm:"type:date"`
	Duration    time.Duration `gorm:"type:integer;not null"`
	Registered  bool          `gorm:"type:integer;default:0"`
}

func (s Song) MarshalJSON() ([]byte, error) {
	type Alias Song

	url := fmt.Sprintf("%s/%s", "https://open.spotify.com/track", s.SpotifyID)
	return json.Marshal(&struct {
		Alias
		Url     string
		Seconds float64
	}{
		Alias:   (Alias)(s),
		Url:     url,
		Seconds: s.Duration.Seconds(),
	})
}

type Share struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	MasterPercent float32 `gorm:"type:real;not null"`
	PubPercent    float32 `gorm:"type:real;not null"`

	UserID uint `gorm:"type:varchar(15);not null;index:idx_user_song,unique"`
	User   User `gorm:"-" json:"-"`
	SongID uint `gorm:"type:varchar(15);not null;index:idx_user_song,unique"`
	Song   Song `gorm:"foreignKey:SongID;references:ID;constraint:OnDelete:CASCADE;"`
}

// type Info struct {
// 	Share Share
// 	Song  Song
// }
