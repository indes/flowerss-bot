package model

type Option struct {
	ID    int64 `gorm:"primary_key"`
	Name  string
	Value string
	EditTime
}
