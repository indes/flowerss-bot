package storage

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/indes/flowerss-bot/internal/model"
)

func GetTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		t.Log(err)
		return nil
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Source{})
	db.AutoMigrate(&model.Subscribe{})
	db.AutoMigrate(&model.Content{})
	return db
}

func TestUserStorageImpl(t *testing.T) {
	db := GetTestDB(t)
	s := NewUserStorageImpl(db)
	ctx := context.Background()

	user := &model.User{
		TelegramID: 123,
		State:      1,
	}

	t.Run(
		"save user", func(t *testing.T) {
			err := s.CrateUser(ctx, user)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"get user", func(t *testing.T) {
			got, err := s.GetUser(ctx, user.ID)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, user.TelegramID, got.TelegramID)
		},
	)

	t.Run(
		"get user by telegram id", func(t *testing.T) {
			got, err := s.GetUserByTelegramID(ctx, user.TelegramID)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, user.TelegramID, got.TelegramID)
			assert.Equal(t, user.ID, got.ID)
		},
	)
}
