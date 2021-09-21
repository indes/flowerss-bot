package model

// Option bot 设置
type Option struct {
	ID    int `gorm:"primary_key;AUTO_INCREMENT"`
	Name  string
	Value string
	EditTime
}
