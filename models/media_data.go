package models

type MediaData struct {
	ID            uint   `gorm:"primaryKey"`
	CarNumber     string `gorm:"size:10"`
	Timestamp     string `gorm:"size:30"`
	PhotoPath1    string `gorm:"size:255"`
	PhotoPath2    string `gorm:"size:255"`
	PhotoPath3    string `gorm:"size:255"`
	VideoPath     string `gorm:"size:255"`
	FirstRequest  bool   `gorm:"default:false"`
	SecondRequest bool   `gorm:"default:false"`
}
