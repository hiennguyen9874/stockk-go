package models

import "github.com/google/uuid"

type Ticker struct {
	Id        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Symbol    string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Exchange  string    `gorm:"type:varchar(100);not null"`
	FullName  string    `gorm:"type:varchar(100);not null"`
	ShortName string    `gorm:"type:varchar(100);not null"`
	Type      string    `gorm:"type:varchar(100);not null"`
	IsActive  bool      `gorm:"not null;default:false"`
}

func (Ticker) TableName() string {
	return "ticker"
}
