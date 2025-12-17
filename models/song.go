package models

import (
	"time"
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
	Title       string        `gorm:"type:varchar(255);not null"`
	Artist      string        `gorm:"type:varchar(255);not null"`
	Iswc        string        `gorm:"type:varchar(15)"`
	Isrc        string        `gorm:"type:varchar(15);uniqueIndex;not null"`
	Upc         uint64        `gorm:""`
	Label       string        `gorm:"type:varchar(255)"`
	ReleaseDate time.Time     `gorm:"type:date"`
	Duration    time.Duration `gorm:"-"`
	Url         string        `gorm:"type:text"`
	Registered  bool          `gorm:"type:integer;default:0"`
}

type Share struct {
	MasterPercent float32
	PubPercent    float32
	Person        Person
}

type Person struct {
	FirstName    string
	LastName     string
	WriterIpiNum uint64
	// PubIpiNum    uint64
	Society Society
	// PubEntity publishers.Entity
}

type Info struct {
	Share Share
	Song  Song
}
