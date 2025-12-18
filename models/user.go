package models

type User struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	Username  string `gorm:"type:varchar(255);uniqueIndex;not null"`
	FirstName string `gorm:"type:varchar(255);not null"`
	LastName  string `gorm:"type:varchar(255);not null"`
	Password  string `gorm:"type:text;not null"`

	Society Society `gorm:"type:text"`

	PubIpi    uint64 `gorm:"type:integer;not null"`
	WriterIpi uint64 `gorm:"type:integer;not null"`
}
