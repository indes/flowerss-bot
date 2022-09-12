package storage

import (
	"context"
	"errors"

	"gorm.io/gorm"

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
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *UserStorageImpl) GetUser(ctx context.Context, id int64) (*model.User, error) {
	var user = &model.User{}
	result := s.db.Where(&model.User{ID: id}).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, result.Error
	}
	return user, nil
}

func (s *UserStorageImpl) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	var user = &model.User{}
	result := s.db.Where(model.User{TelegramID: telegramID}).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, result.Error
	}
	return user, nil
}