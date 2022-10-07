package model

// Content fetcher content
type Content struct {
	SourceID     uint
	HashID       string `gorm:"primary_key"`
	RawID        string
	RawLink      string
	Title        string
	Description  string `gorm:"-"` //ignore to db
	TelegraphURL string
	EditTime
}
