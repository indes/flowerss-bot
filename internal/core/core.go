package core

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/internal/storage"
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
		zap.S().Fatalf("connect db failed, err: %+v", err)
		return nil
	}
	return &Core{
		userStorage:         storage.NewUserStorageImpl(db),
		contentStorage:      storage.NewContentStorageImpl(db),
		sourceStorage:       storage.NewSourceStorageImpl(db),
		subscriptionStorage: storage.NewSubscriptionStorageImpl(db),
	}
}

func (c *Core) Run() error {
	go func() {
		zap.S().Infoln("core running!")
		ctx := context.Background()
		for true {
			count, err := c.subscriptionStorage.CountSubscriptions(ctx)
			if err != nil {
				zap.S().Errorln(err)
			} else {
				zap.S().Infof("CountSubscriptions %v", count)
			}
			time.Sleep(3 * time.Second)
		}
	}()
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
