package models

type DataLink struct {
	LongUrl   string
	Teacher   string
	ClassInfo string
	Token     string
	Lesson    string
}

type ShortLink struct {
	ID       int
	ShortUrl string
	Data     *DataLink
}

type Links struct {
	DataLinks []*DataLink
}
