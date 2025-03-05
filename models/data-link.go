package models

import "gorm.io/gorm"

type DataLink struct {
	gorm.Model
	LongUrl   string `json:"long_url"`
	Teacher   string `json:"teacher"`
	ClassInfo string `json:"class_info"`
	Token     string `json:"token"`
	Lesson    string `json:"lesson"`
}

type ShortLink struct {
	gorm.Model
	ShortUrl   string   `json:"short_url"`
	UniqueCode string   `json:"unique_code"`
	DataLinkID uint     `json:"data_id"`
	DataLink   DataLink `json:"data_link"   gorm:"foreignKey:DataLinkID"`
	QrImageUrl string   `json:"qrimage_url"`
}

type Links struct {
	DataLinks []*DataLink `json:"data_links"`
}
