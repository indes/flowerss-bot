package model

type Option struct {
	ID    int `gorm:"primary_key;AUTO_INCREMENT"`
	Name  string
	Value string
	EditTime
}
