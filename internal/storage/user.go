package storage

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/indes/flowerss-bot/internal/model"
)

type UserStorageImpl struct {
	db *gorm.DB
}

func NewUserStorageImpl(db *gorm.DB) *UserStorageImpl {
	return &UserStorageImpl{db: db}
}

func (s *UserStorageImpl) CrateUser(ctx context.Context, user *model.User) error {
	result := s.db.Create(user)
	if len(result.GetErrors()) != 0 {
		return result.GetErrors()[0]
	}
	return nil
}

func (s *UserStorageImpl) GetUser(ctx context.Context, id int64) (*model.User, error) {
	var user = &model.User{}
	db := s.db.Where(&model.User{ID: id}).First(user)
	if db.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	if len(db.GetErrors()) != 0 {
		return nil, db.GetErrors()[0]
	}
	return user, nil
}

func (s *UserStorageImpl) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	var user = &model.User{}
	db := s.db.Where(model.User{TelegramID: telegramID}).First(user)
	if db.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	if len(db.GetErrors()) != 0 {
		return nil, db.GetErrors()[0]
	}
	return user, nil
}
