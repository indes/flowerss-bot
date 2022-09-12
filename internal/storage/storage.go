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

type Storage interface {
	Init(ctx context.Context) error
}

// UserStorage 用户存储接口
type UserStorage interface {
	Storage
	CrateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
}

// SourceStorage 订阅源存储接口
type SourceStorage interface {
	Storage
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
	Count    int // 需要获取的数量，-1为获取全部
	Offset   int
	SortType SubscriptionSortType
}

type GetSubscriptionsResult struct {
	Subscriptions []*model.Subscribe
	HasMore       bool
}

type SubscriptionStorage interface {
	Storage
	AddSubscription(ctx context.Context, subscription *model.Subscribe) error
	GetSubscriptionsByUserID(
		ctx context.Context, userID int64, opts *GetSubscriptionsOptions,
	) (*GetSubscriptionsResult, error)
	GetSubscriptionsBySourceID(
		ctx context.Context, sourceID uint, opts *GetSubscriptionsOptions,
	) (*GetSubscriptionsResult, error)
	CountSubscriptions(ctx context.Context) (int64, error)
	DeleteSubscription(ctx context.Context, userID int64, sourceID uint) (int64, error)
}

type ContentStorage interface {
	Storage
	// AddContent 添加一条文章
	AddContent(ctx context.Context, content *model.Content) error
	// DeleteSourceContents 删除订阅源的所有文章，返回被删除的文章数
	DeleteSourceContents(ctx context.Context, sourceID uint) (int64, error)
	// HashIDExist hash id 对应的文章是否已存在
	HashIDExist(ctx context.Context, hashID string) (bool, error)
}
