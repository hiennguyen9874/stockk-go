package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type WatchList struct {
	gorm.Model
	Id        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	Name      string         `gorm:"type:varchar(100);not null"`
	Tickers   pq.StringArray `gorm:"type:text[]"`
	OwnerId   uint
}

func (WatchList) TableName() string {
	return "watchlist"
}
