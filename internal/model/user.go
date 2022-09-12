package model

// User subscriber
//
// TelegramID 用作外键
type User struct {
	ID         int64 `gorm:"primary_key"`
	TelegramID int64
	State      int `gorm:"DEFAULT:0;"`
	EditTime
}

// FindOrCreateUserByTelegramID find subscriber or init a subscriber by telegram ID
func FindOrCreateUserByTelegramID(telegramID int64) (*User, error) {
	var user User
	db.Where(User{TelegramID: telegramID}).FirstOrCreate(&user)

	return &user, nil
}
