package models

import (
	"time"
)

type Chart struct {
	Id           uint      `gorm:"primaryKey"`
	OwnerSource  string    `gorm:"type:varchar(100);index;not null"`
	OwnerId      string    `gorm:"type:varchar(100);index;not null"`
	Name         string    `gorm:"type:varchar(100);not null"`
	Symbol       string    `gorm:"type:varchar(100);not null"`
	Resolution   string    `gorm:"type:varchar(100);not null"`
	LastModified time.Time `gorm:"autoUpdateTime;not null"`
	Content      string    `gorm:"type:text;not null"`
}
