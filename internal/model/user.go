package model

// User subscriber
//
// TelegramID 用作外键
type User struct {
	ID int64 `gorm:"primary_key"`
	EditTime
}
