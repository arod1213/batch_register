package publishers

import "errors"

type Entity string

const (
	UniversalMusicPublishing Entity = "Universal Music Publishing Group"

	SonyATVTunes  Entity = "SONY/ATV TUNES LLC"
	SonyATVSongs  Entity = "SONY/ATV SONGS LLC"
	SonyATVSONATA Entity = "SONY/ATV SONATA"

	WCMusicCorp               Entity = "WC MUSIC CORP."
	WarnerTamerlanePublishing Entity = "WARNER-TAMERLANE PUBLISHING CORP."

	ReservoirMusic Entity = "Reservoir Music"
	Songtrust      Entity = "Songtrust"
	Kobalt         Entity = "Kobalt Music Group"
)

type Publisher struct {
	Name   string
	IpiNum uint64
	Number string
	Entity Entity
}

func (e Entity) AdminOrPub() Publisher {
	var name string
	var ipiNum uint64
	var number string

	switch e {
	case ReservoirMusic:
		name = "Reservoir Music"
		ipiNum = 551650169
		number = "P41349"

	case SonyATVSONATA:
		name = "SONY/ATV SONATA"
		number = "P90541"
		ipiNum = 689749753

	case SonyATVSongs:
		name = "SONY/ATV SONGS LLC"
		number = "P8301A"
		ipiNum = 187062752

	case SonyATVTunes:
		name = "SONY/ATV TUNES LLC"
		ipiNum = 338164558
		number = "P8301C"

	case UniversalMusicPublishing:
		name = "SONGS OF UNIVERSAL, INC."
		ipiNum = 353271280
		number = "P1195V"

	case WCMusicCorp:
		name = "WC MUSIC CORP."
		ipiNum = 392888203
		number = "P93725"

	case WarnerTamerlanePublishing:
		name = "WARNER-TAMERLANE PUBLISHING CORP."
		ipiNum = 185314175
		number = "P94075"

	case Songtrust:
		name = "Songtrust"
		ipiNum = 613926842
		number = "P8368C"

	case Kobalt:
		name = "KOBALT MUSIC PUB AMERICA INC"
		ipiNum = 503659557
		number = "P4614Q"

	default:
		panic("unknown entity: " + string(e))
	}

	return Publisher{
		Name:   name,
		Number: number,
		IpiNum: ipiNum,
		Entity: e,
	}
}

func (e Entity) Publisher() (*Publisher, error) {
	var name string
	var ipiNum uint64
	var number string

	switch e {
	case Kobalt:
		name = "ARTIST PUB GROUP WEST"
		ipiNum = 480980625
		number = "P8446Q"
	default:
		return nil, errors.New("unknown entity: " + string(e))
	}

	pub := &Publisher{
		Name:   name,
		Number: number,
		IpiNum: ipiNum,
		Entity: e,
	}
	return pub, nil
}
