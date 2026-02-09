package models

type Checklist struct {
	UserID uint `gorm:"type:varchar(15);not null;index:idx_user_song,unique"`
	User   User `gorm:"-" json:"-"`
	SongID uint `gorm:"type:varchar(15);not null;index:idx_user_song,unique"`
	Song   Song `gorm:"foreignKey:SongID;references:ID;constraint:OnDelete:CASCADE;"`
}
