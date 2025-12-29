package models

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	Username  string `gorm:"type:varchar(255);uniqueIndex;not null"`
	FirstName string `gorm:"type:varchar(255);not null"`
	LastName  string `gorm:"type:varchar(255);not null"`
	Password  string `gorm:"type:text;not null"`

	Society Society `gorm:"type:text"`

	PubIpi    uint64  `gorm:"type:integer;not null"`
	WriterIpi uint64  `gorm:"type:integer;not null"`
	DiscogID  *string `gorm:"type:text"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User

	var url *string = nil
	if u.DiscogID != nil {
		x := fmt.Sprintf("%s/%s", "https://open.spotify.com/playlist", *u.DiscogID)
		url = &x
	}

	return json.Marshal(&struct {
		Alias
		DiscogUrl *string
	}{
		Alias:     (Alias)(u),
		DiscogUrl: url,
	})
}
