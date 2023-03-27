package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Id                uint      `gorm:"primaryKey"`
	CreatedAt         time.Time `gorm:"not null;default:now()"`
	UpdatedAt         time.Time `gorm:"not null;default:now()"`
	CurrentTicker     *string   `gorm:"type:varchar(100)"`
	CurrentResolution *string   `gorm:"type:varchar(100)"`
	OwnerId           uint
}

func (Client) TableName() string {
	return "client"
}
