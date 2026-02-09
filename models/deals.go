package models

// type Deal struct {
// 	SongID uint `gorm:"type:varchar(15);not null"`
// 	Song   Song `gorm:"foreignKey:SongID;references:ID;constraint:OnDelete:CASCADE;"`
//
// 	UserID uint `gorm:"type:varchar(15);not null"`
// 	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
// }

type DealType interface {
	GetSongID() uint
	SetSongID(uint)
}

type Credit struct {
	SongID uint `gorm:"type:varchar(15);not null"`
	Song   Song `gorm:"foreignKey:SongID;references:ID;constraint:OnDelete:CASCADE;"`

	// use name if not linking to verified user
	Name string `gorm:"type:text;not null"`

	UserID *uint `gorm:"type:varchar(15)"`
	User   *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	CreatedByUserID uint `gorm:"type:varchar(15);not null"`
	CreatedByUser   User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	Role string `gorm:"type:text;not null"`
}

type PubDeal struct {
	OverridingID   *uint    `goSongIDrm:"type:varchar(15)"`
	OverridingDeal *PubDeal `gorm:"foreignKey:OverridingID;references:ID;constraint:OnDelete:CASCADE;"`

	SongID uint `gorm:"type:varchar(15);not null"`
	Song   Song `gorm:"foreignKey:SongID;references:ID;constraint:OnDelete:CASCADE;"`

	// use name if not linking to verified user
	Name      string  `gorm:"type:text;not null"`
	AdminInfo *string `gorm:"type:text"`

	UserID *uint `gorm:"type:varchar(15)"`
	User   *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	CreatedByUserID uint `gorm:"type:varchar(15);not null"`
	CreatedByUser   User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	Percent float32 `gorm:"type:real;not null"`
}

func (d PubDeal) GetSongID() uint {
	return d.SongID
}
func (d *PubDeal) SetSongID(id uint) {
	d.SongID = id
}

type RoyaltyType string

const (
	PPD   RoyaltyType = "PPD"
	Net   RoyaltyType = "NetReceipts"
	Gross RoyaltyType = "Gross"
)

type MasterDeal struct {
	OverridingID   *uint       `goSongIDrm:"type:varchar(15)"`
	OverridingDeal *MasterDeal `gorm:"foreignKey:OverridingID;references:ID;constraint:OnDelete:CASCADE;"`

	SongID uint `gorm:"type:varchar(15);not null"`
	Song   Song `gorm:"foreignKey:SongID;references:ID;constraint:OnDelete:CASCADE;"`

	// use name if not linking to verified user
	Name string `gorm:"type:text;not null"`

	UserID *uint `gorm:"type:varchar(15)"`
	User   *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	CreatedByUserID uint `gorm:"type:varchar(15);not null"`
	CreatedByUser   User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`

	Fee          float32     `gorm:"type:real;not null"`
	Advance      float32     `gorm:"type:real;not null"`
	RecoupAmount float32     `gorm:"type:real;not null"`
	Points       float32     `gorm:"type:real;not null"`
	RoyaltyType  RoyaltyType `gorm:"type:text;not null"`
	Territory    string      `gorm:"type:text;not null"`
}

func (d MasterDeal) GetSongID() uint {
	return d.SongID
}

func (d *MasterDeal) SetSongID(id uint) {
	d.SongID = id
}
