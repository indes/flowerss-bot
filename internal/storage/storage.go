package storage

import (
	"context"
	"errors"

	"github.com/indes/flowerss-bot/internal/model"
)

var (
	// ErrRecordNotFound 数据不存在错误
	ErrRecordNotFound = errors.New("record not found")
)

// UserStorage 用户存储接口
type UserStorage interface {
	CrateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
}

// SourceStorage 订阅源存储接口
type SourceStorage interface {
	AddSource(ctx context.Context, source *model.Source) error
	GetSource(ctx context.Context, sourceID uint) (*model.Source, error)
	GetSourceByURL(ctx context.Context, url string) (*model.Source, error)
}

type FeedStorage interface {
}

type SubscriptionSortType = int

const (
	SubscriptionSortTypeCreatedTimeDesc SubscriptionSortType = iota
)

type GetSubscriptionsOptions struct {
	Count      int // 需要获取的数量，-1为获取全部
	Offset     int
	AttachInfo string
	SortType   SubscriptionSortType
}

type GetSubscriptionsResult struct {
	Subscriptions []*model.Subscribe
	HasMore       bool
	AttachInfo    string
}

type SubscriptionStorage interface {
	AddSubscription(ctx context.Context, subscription *model.Subscribe) error
	GetSubscriptionsByUserID(
		ctx context.Context, userID int64, opts *GetSubscriptionsOptions,
	) (*GetSubscriptionsResult, error)
	GetSubscriptionsBySourceID(
		ctx context.Context, sourceID uint, opts *GetSubscriptionsOptions,
	) (*GetSubscriptionsResult, error)
}
