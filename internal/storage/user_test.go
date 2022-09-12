package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/indes/flowerss-bot/internal/model"
)

func GetTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		t.Log(err)
		return nil
	}
	return db.Debug()
}

func TestUserStorageImpl(t *testing.T) {
	db := GetTestDB(t)
	s := NewUserStorageImpl(db)
	ctx := context.Background()
	s.Init(ctx)
	user := &model.User{
		ID: 123,
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
			assert.Equal(t, user.ID, got.ID)
		},
	)
}
