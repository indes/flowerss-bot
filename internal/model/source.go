package model

type Source struct {
	ID         uint `gorm:"primary_key;AUTO_INCREMENT"`
	Link       string
	Title      string
	ErrorCount uint
	Content    []Content
	EditTime
}
