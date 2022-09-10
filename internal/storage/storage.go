package storage

import (
	"context"
	"errors"

	"github.com/indes/flowerss-bot/internal/model"
)

var (
	// ErrRecordNotFound returns a "record not found error".
	ErrRecordNotFound = errors.New("record not found")
)

type UserStorage interface {
	CrateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
}

type FeedStorage interface {
}

type SourceStorage interface {
}

type SubscriptionStorage interface {
}
