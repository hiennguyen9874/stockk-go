package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Email       string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password    string    `gorm:"type:varchar(100);not null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	IsActive    bool      `gorm:"not null"`
	IsSuperUser bool      `gorm:"not null"`
}
