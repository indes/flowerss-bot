package model

import "errors"

// User subscriber
//
// TelegramID 用作外键
type User struct {
	ID         int64 `gorm:"primary_key"`
	TelegramID int64
	Source     []Source `gorm:"many2many:subscribes;"`
	State      int      `gorm:"DEFAULT:0;"`
	EditTime
}

// FindOrInitUser find subscriber or init a subscriber if user not in db
//
// Deprecated: Use model.FindOrCreateUserByTelegramID instead.
func FindOrInitUser(userID int64) (*User, error) {
	var user User
	db.Where(User{ID: userID}).FirstOrCreate(&user)

	return &user, nil
}

// FindOrCreateUserByTelegramID find subscriber or init a subscriber by telegram ID
func FindOrCreateUserByTelegramID(telegramID int64) (*User, error) {
	var user User
	db.Where(User{TelegramID: telegramID}).FirstOrCreate(&user)

	return &user, nil
}

// GetSubSourceMap get user subscribe and fetcher source
func (user *User) GetSubSourceMap() (map[Subscribe]Source, error) {
	m := make(map[Subscribe]Source)

	subs, err := GetSubsByUserID(user.TelegramID)
	if err != nil {
		return nil, errors.New("get subs error")
	}

	for _, sub := range subs {
		source, err := GetSourceById(sub.SourceID)
		if err != nil {
			continue
		}
		m[sub] = *source
	}

	return m, nil
}

// GetSubSourceList get user subscribe and fetcher source
func (user *User) GetSubSourceList() ([]*SubWithSource, error) {
	m := make(map[Subscribe]struct{})

	subs, err := GetSubsByUserID(user.TelegramID)
	if err != nil {
		return nil, errors.New("get subs error")
	}

	subsWithSources := make([]*SubWithSource, 0, len(subs))
	for _, sub := range subs {
		sub := sub
		if _, ok := m[sub]; ok {
			continue
		}
		m[sub] = struct{}{}
		source, err := GetSourceById(sub.SourceID)
		if err != nil {
			continue
		}
		subsWithSources = append(subsWithSources, &SubWithSource{
			Sub: &sub,
			Src: source,
		})
	}

	return subsWithSources, nil
}
