package core

import (
	"context"
	"errors"
	"strings"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/indes/flowerss-bot/internal/config"
	"github.com/indes/flowerss-bot/internal/core/task"
	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/internal/storage"
)

var (
	ErrSubscriptionExist    = errors.New("already subscribed")
	ErrSubscriptionNotExist = errors.New("subscription not exist")
	ErrSourceNotExist       = errors.New("source not exist")
)

type Core struct {
	// Storage
	userStorage         storage.User
	contentStorage      storage.Content
	sourceStorage       storage.Source
	subscriptionStorage storage.Subscription

	rssTask *task.RssUpdateTask
}

func NewCore(
	userStorage storage.User, contentStorage storage.Content, sourceStorage storage.Source,
	subscriptionStorage storage.Subscription, rssTask *task.RssUpdateTask,
) *Core {
	return &Core{
		userStorage:         userStorage,
		contentStorage:      contentStorage,
		sourceStorage:       sourceStorage,
		subscriptionStorage: subscriptionStorage,
		rssTask:             rssTask,
	}
}

func NewCoreFormConfig() *Core {
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

	return NewCore(
		storage.NewUserStorageImpl(db),
		storage.NewContentStorageImpl(db),
		storage.NewSourceStorageImpl(db),
		storage.NewSubscriptionStorageImpl(db),
		task.NewRssTask(),
	)
}

func (c *Core) Init() error {
	c.userStorage.Init(context.Background())
	c.contentStorage.Init(context.Background())
	c.sourceStorage.Init(context.Background())
	c.subscriptionStorage.Init(context.Background())
	return nil
}

func (c *Core) Run() error {
	c.rssTask.Start()
	return nil
}

func (c *Core) RegisterRssUpdateObserver(o task.RssUpdateObserver) {
	c.rssTask.Register(o)
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
			log.Errorf("get source %d failed, %v", subs.SourceID, err)
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

// Unsubscribe 添加订阅
func (c *Core) Unsubscribe(ctx context.Context, userID int64, sourceID uint) error {
	exist, err := c.subscriptionStorage.SubscriptionExist(ctx, userID, sourceID)
	if err != nil {
		return err
	}

	if !exist {
		return ErrSubscriptionNotExist
	}

	// 移除该用户订阅
	_, err = c.subscriptionStorage.DeleteSubscription(ctx, userID, sourceID)
	if err != nil {
		return err
	}

	// 获取源的订阅数量
	count, err := c.subscriptionStorage.CountSourceSubscriptions(ctx, sourceID)
	if err != nil {
		return err
	}

	if count != 0 {
		return nil
	}

	// 如果源不再有订阅用户，移除该订阅源
	if err := c.removeSource(ctx, sourceID); err != nil {
		return err
	}
	return nil
}

// removeSource 移除订阅源
func (c *Core) removeSource(ctx context.Context, sourceID uint) error {
	if err := c.sourceStorage.Delete(ctx, sourceID); err != nil {
		return err
	}

	count, err := c.contentStorage.DeleteSourceContents(ctx, sourceID)
	if err != nil {
		return err
	}
	log.Infof("remove source %d and %d contents", sourceID, count)
	return nil
}

// GetSourceByURL 获取用户订阅的订阅源
func (c *Core) GetSourceByURL(ctx context.Context, sourceURL string) (*model.Source, error) {
	source, err := c.sourceStorage.GetSourceByURL(ctx, sourceURL)
	if err != nil {
		if err == storage.ErrRecordNotFound {
			return nil, ErrSourceNotExist
		}
		return nil, err
	}
	return source, nil
}

// GetSource 获取用户订阅的订阅源
func (c *Core) GetSource(ctx context.Context, id uint) (*model.Source, error) {
	source, err := c.sourceStorage.GetSource(ctx, id)
	if err != nil {
		if err == storage.ErrRecordNotFound {
			return nil, ErrSourceNotExist
		}
		return nil, err
	}
	return source, nil
}

// UnsubscribeAllSource 添加订阅
func (c *Core) UnsubscribeAllSource(ctx context.Context, userID int64) error {
	sources, err := c.GetUserSubscribedSources(ctx, userID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := range sources {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			if err := c.Unsubscribe(ctx, userID, sources[i].ID); err != nil {
				log.Errorf("user %d unsubscribe %d failed, %v", userID, sources[i].ID, err)
			}
		}()
	}
	wg.Wait()
	return nil
}

// GetSubscription 获取订阅
func (c *Core) GetSubscription(ctx context.Context, userID int64, sourceID uint) (*model.Subscribe, error) {
	subscription, err := c.subscriptionStorage.GetSubscription(ctx, userID, sourceID)
	if err != nil {
		if err == storage.ErrRecordNotFound {
			return nil, ErrSubscriptionNotExist
		}
		return nil, err
	}
	return subscription, nil
}

// SetSubscriptionTag 设置订阅标签
func (c *Core) SetSubscriptionTag(ctx context.Context, userID int64, sourceID uint, tags []string) error {
	subscription, err := c.GetSubscription(ctx, userID, sourceID)
	if err != nil {
		return err
	}

	subscription.Tag = "#" + strings.Join(tags, " #")
	return c.subscriptionStorage.UpdateSubscription(ctx, userID, sourceID, subscription)
}

// SetSubscriptionInterval
func (c *Core) SetSubscriptionInterval(ctx context.Context, userID int64, sourceID uint, interval int) error {
	subscription, err := c.GetSubscription(ctx, userID, sourceID)
	if err != nil {
		return err
	}

	subscription.Interval = interval
	return c.subscriptionStorage.UpdateSubscription(ctx, userID, sourceID, subscription)
}

// EnableSourceUpdate 开启订阅源更新
func (c *Core) EnableSourceUpdate(ctx context.Context, sourceID uint) error {
	return c.ClearSourceErrorCount(ctx, sourceID)
}

// DisableSourceUpdate 关闭订阅源更新
func (c *Core) DisableSourceUpdate(ctx context.Context, sourceID uint) error {
	source, err := c.GetSource(ctx, sourceID)
	if err != nil {
		return err
	}

	source.ErrorCount = config.ErrorThreshold + 1
	return c.sourceStorage.UpsertSource(ctx, sourceID, source)
}

// ClearSourceErrorCount 清空订阅源错误计数
func (c *Core) ClearSourceErrorCount(ctx context.Context, sourceID uint) error {
	source, err := c.GetSource(ctx, sourceID)
	if err != nil {
		return err
	}

	source.ErrorCount = 0
	return c.sourceStorage.UpsertSource(ctx, sourceID, source)
}

func (c *Core) ToggleSubscriptionNotice(ctx context.Context, userID int64, sourceID uint) error {
	subscription, err := c.GetSubscription(ctx, userID, sourceID)
	if err != nil {
		return err
	}
	if subscription.EnableNotification == 1 {
		subscription.EnableNotification = 0
	} else {
		subscription.EnableNotification = 1
	}
	return c.subscriptionStorage.UpsertSubscription(ctx, userID, sourceID, subscription)
}

func (c *Core) ToggleSourceUpdateStatus(ctx context.Context, sourceID uint) error {
	source, err := c.GetSource(ctx, sourceID)
	if err != nil {
		return err
	}

	if source.ErrorCount < config.ErrorThreshold {
		source.ErrorCount = config.ErrorThreshold + 1
	} else {
		source.ErrorCount = 0
	}
	return c.sourceStorage.UpsertSource(ctx, sourceID, source)
}

func (c *Core) ToggleSubscriptionTelegraph(ctx context.Context, userID int64, sourceID uint) error {
	subscription, err := c.GetSubscription(ctx, userID, sourceID)
	if err != nil {
		return err
	}
	if subscription.EnableTelegraph == 1 {
		subscription.EnableTelegraph = 0
	} else {
		subscription.EnableTelegraph = 1
	}
	return c.subscriptionStorage.UpsertSubscription(ctx, userID, sourceID, subscription)
}
