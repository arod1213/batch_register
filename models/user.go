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
	GeniusID  *uint   `gorm:"type:integer"`

	// Workspace onboarding / preferences (server-backed; no local storage).
	WorkspaceName        string `gorm:"type:varchar(255);not null;default:''"`
	WorkspaceRole        string `gorm:"type:text;not null;default:''"`       // ARTIST|MANAGER|PUBLISHER|LABEL|OPS
	WorkspacePrimaryUse  string `gorm:"type:text;not null;default:'BOTH'"`   // REGISTRATIONS|ROYALTIES|BOTH
	WorkspaceCreatorType string `gorm:"type:text;not null;default:'ARTIST'"` // ARTIST|WRITER|PRODUCER (when role=ARTIST)
	WorkspaceOnboarded   bool   `gorm:"type:boolean;not null;default:false"`
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
