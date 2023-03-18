package models

type Ticker struct {
	Id        uint   `gorm:"primaryKey"`
	Symbol    string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Exchange  string `gorm:"type:varchar(100);not null"`
	FullName  string `gorm:"type:varchar(100);not null"`
	ShortName string `gorm:"type:varchar(100);not null"`
	Type      string `gorm:"type:varchar(100);not null"`
	IsActive  bool   `gorm:"not null;default:false"`
}

func (Ticker) TableName() string {
	return "ticker"
}
