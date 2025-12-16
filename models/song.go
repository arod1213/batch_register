package models

import "time"

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
		// case HFA:
		// 	return "034"
	}
	panic("unreachable - fell through switch")
}

type Song struct {
	Title       string
	Artist      string
	Iswc        string
	Isrc        string
	Upc         uint64
	Label       string
	ReleaseDate time.Time
	Duration    time.Duration
	Url         string
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
	PubIpiNum    uint64
	Society      Society
}

type Info struct {
	Share Share
	Song  Song
}
