package model

// User subscriber
type User struct {
	ID int64 `gorm:"primary_key"`
	EditTime
}
