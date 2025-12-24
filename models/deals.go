package models

type Deal struct {
	SongID uint `gorm:"type:varchar(15);not null"`
	Song   Song `gorm:"foreignKey:SongID;references:ID;constraint:OnDelete:CASCADE;"`

	UserID uint `gorm:"type:varchar(15);not null"`
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
}
