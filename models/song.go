package models

import (
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
	Title  string `gorm:"type:varchar(255);not null"`
	Artist string `gorm:"type:varchar(255);not null"`
	Label  string `gorm:"type:varchar(255)"`

	Iswc *string `gorm:"type:varchar(15)"`
	Isrc string  `gorm:"primaryKey;type:varchar(15)"`
	Upc  uint64  `gorm:"type:integer;not null"`

	Url         string        `gorm:"type:text"`
	ReleaseDate time.Time     `gorm:"type:date"`
	Duration    time.Duration `gorm:"-"`
	Registered  bool          `gorm:"type:integer;default:0"`

	Share *Share `gorm:"foreignKey:SongIsrc;references:Isrc;->"`
}

func (s *Song) AfterFind(tx *gorm.DB) error {
	if s.Share == nil {
		share := Share{}
		s.Share = &share
	}
	return nil
}

type Share struct {
	MasterPercent float32 `gorm:"type:real;not null"`
	PubPercent    float32 `gorm:"type:real;not null"`
	Person        Person  `gorm:"-" json:"-"`

	SongIsrc string `gorm:"primaryKey;type:varchar(15)" json:"-"`
	Song     Song   `gorm:"foreignKey:SongIsrc;references:Isrc;constraint:OnDelete:CASCADE;" json:"-"`
}

type Person struct {
	FirstName    string
	LastName     string
	WriterIpiNum uint64
	// PubIpiNum    uint64
	Society Society
	// PubEntity publishers.Entity
}

// type Info struct {
// 	Share Share
// 	Song  Song
// }
