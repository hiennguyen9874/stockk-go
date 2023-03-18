package models

type StudyTemplate struct {
	Id          uint   `gorm:"primaryKey"`
	OwnerSource string `gorm:"type:varchar(100);index;not null"`
	OwnerId     string `gorm:"type:varchar(100);index;not null"`
	Name        string `gorm:"type:varchar(100);not null"`
	Content     string `gorm:"type:text;not null"`
}
