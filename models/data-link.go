package models

type DataLink struct {
	LongUrl   string `json:"long_url"`
	Teacher   string `json:"teacher"`
	ClassInfo string `json:"class_info"`
	Token     string `json:"token"`
	Lesson    string `json:"lesson"`
}

type ShortLink struct {
	ID         int       `json:"id"`
	ShortUrl   string    `json:"short_url"`
	UniqueCode string    `json:"unique_code"`
	Data       *DataLink `json:"data"`
	QrImageUrl string    `json:"qrimage_url"`
}

type Links struct {
	DataLinks []*DataLink
}
