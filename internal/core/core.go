package core

import (
	"context"
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/internal/storage"
)

var (
	ErrSubscriptionExist = errors.New("already subscribed")
)

type Core struct {
	// Storage
	userStorage         storage.User
	contentStorage      storage.Content
	sourceStorage       storage.Source
	subscriptionStorage storage.Subscription
}

func NewCore() *Core {
	var err error
	var db *gorm.DB
	if config.EnableMysql {
		db, err = gorm.Open(mysql.Open(config.Mysql.GetMysqlConnectingString()))
	} else {
		db, err = gorm.Open(sqlite.Open(config.SQLitePath))
	}
	if err != nil {
		log.Fatalf("connect db failed, err: %+v", err)
		return nil
	}

	if config.DBLogMode {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)

	return &Core{
		userStorage:         storage.NewUserStorageImpl(db),
		contentStorage:      storage.NewContentStorageImpl(db),
		sourceStorage:       storage.NewSourceStorageImpl(db),
		subscriptionStorage: storage.NewSubscriptionStorageImpl(db),
	}
}

func (c *Core) Init() error {
	c.userStorage.Init(context.Background())
	c.contentStorage.Init(context.Background())
	c.sourceStorage.Init(context.Background())
	c.subscriptionStorage.Init(context.Background())
	return nil
}

// GetUserSubscribedSources 获取用户订阅的订阅源
func (c *Core) GetUserSubscribedSources(ctx context.Context, userID int64) ([]*model.Source, error) {
	opt := &storage.GetSubscriptionsOptions{Count: -1}
	result, err := c.subscriptionStorage.GetSubscriptionsByUserID(ctx, userID, opt)
	if err != nil {
		return nil, err
	}

	var sources []*model.Source
	for _, subs := range result.Subscriptions {
		source, err := c.sourceStorage.GetSource(ctx, subs.SourceID)
		if err != nil {
			log.Errorf("get source failed, %v", err)
			continue
		}
		sources = append(sources, source)
	}
	return sources, nil
}

// AddSubscription 添加订阅
func (c *Core) AddSubscription(ctx context.Context, userID int64, sourceID uint) error {
	exist, err := c.subscriptionStorage.SubscriptionExist(ctx, userID, sourceID)
	if err != nil {
		return err
	}

	if exist {
		return ErrSubscriptionExist
	}

	subscription := &model.Subscribe{
		UserID:             userID,
		SourceID:           sourceID,
		EnableNotification: 1,
		EnableTelegraph:    1,
		Interval:           config.UpdateInterval,
		WaitTime:           config.UpdateInterval,
	}
	return c.subscriptionStorage.AddSubscription(ctx, subscription)
}
