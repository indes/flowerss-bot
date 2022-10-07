package storage

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/indes/flowerss-bot/internal/log"
	"github.com/indes/flowerss-bot/internal/model"
)

type SubscriptionStorageImpl struct {
	db *gorm.DB
}

func NewSubscriptionStorageImpl(db *gorm.DB) *SubscriptionStorageImpl {
	return &SubscriptionStorageImpl{db: db.Model(&model.Subscribe{})}
}

func (s *SubscriptionStorageImpl) Init(ctx context.Context) error {
	return s.db.Migrator().AutoMigrate(&model.Subscribe{})
}

func (s *SubscriptionStorageImpl) AddSubscription(ctx context.Context, subscription *model.Subscribe) error {
	result := s.db.WithContext(ctx).Create(subscription)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SubscriptionStorageImpl) SubscriptionExist(ctx context.Context, userID int64, sourceID uint) (bool, error) {
	var count int64
	result := s.db.WithContext(ctx).Where("user_id = ? and source_id = ?", userID, sourceID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return (count > 0), nil
}

func (s *SubscriptionStorageImpl) GetSubscription(ctx context.Context, userID int64, sourceID uint) (
	*model.Subscribe, error,
) {
	subscription := &model.Subscribe{}
	result := s.db.WithContext(ctx).Where("user_id = ? and source_id = ?", userID, sourceID).First(subscription)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, result.Error
	}
	return subscription, nil
}

func (s *SubscriptionStorageImpl) GetSubscriptionsByUserID(
	ctx context.Context, userID int64, opts *GetSubscriptionsOptions,
) (*GetSubscriptionsResult, error) {
	var subscriptions []*model.Subscribe

	count := s.getSubscriptionsCount(opts)
	orderBy := s.getSubscriptionsOrderBy(opts)
	dbResult := s.db.WithContext(ctx).Where(
		&model.Subscribe{UserID: userID},
	).Limit(count).Order(orderBy).Offset(opts.Offset).Find(&subscriptions)
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, dbResult.Error
	}

	result := &GetSubscriptionsResult{}
	if opts.Count != -1 && len(subscriptions) > opts.Count {
		result.HasMore = true
		subscriptions = subscriptions[:opts.Count]
	}

	result.Subscriptions = subscriptions
	return result, nil
}

func (s *SubscriptionStorageImpl) GetSubscriptionsBySourceID(
	ctx context.Context, sourceID uint, opts *GetSubscriptionsOptions,
) (*GetSubscriptionsResult, error) {
	var subscriptions []*model.Subscribe

	count := s.getSubscriptionsCount(opts)
	orderBy := s.getSubscriptionsOrderBy(opts)
	dbResult := s.db.WithContext(ctx).Where(
		&model.Subscribe{SourceID: sourceID},
	).Limit(count).Order(orderBy).Offset(opts.Offset).Find(&subscriptions)
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, dbResult.Error
	}

	result := &GetSubscriptionsResult{}
	if opts.Count > 0 && len(subscriptions) > opts.Count {
		result.HasMore = true
		subscriptions = subscriptions[:opts.Count]
	}

	result.Subscriptions = subscriptions
	return result, nil
}

func (s *SubscriptionStorageImpl) getSubscriptionsCount(opts *GetSubscriptionsOptions) int {
	count := opts.Count
	if count != -1 {
		count += 1
	}
	return count
}

func (s *SubscriptionStorageImpl) getSubscriptionsOrderBy(opts *GetSubscriptionsOptions) string {
	switch opts.SortType {
	case SubscriptionSortTypeCreatedTimeDesc:
		return "created_at desc"
	}
	return ""
}

func (s *SubscriptionStorageImpl) CountSubscriptions(ctx context.Context) (int64, error) {
	var count int64
	dbResult := s.db.WithContext(ctx).Count(&count)
	if dbResult.Error != nil {
		return 0, dbResult.Error
	}
	return count, nil
}

func (s *SubscriptionStorageImpl) DeleteSubscription(ctx context.Context, userID int64, sourceID uint) (int64, error) {
	result := s.db.WithContext(ctx).Where(
		"user_id = ? and source_id = ?", userID, sourceID,
	).Delete(&model.Subscribe{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (s *SubscriptionStorageImpl) CountSourceSubscriptions(ctx context.Context, sourceID uint) (int64, error) {
	var count int64
	result := s.db.WithContext(ctx).Where("source_id = ?", sourceID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (s *SubscriptionStorageImpl) UpdateSubscription(
	ctx context.Context, userID int64, sourceID uint, newSubscription *model.Subscribe,
) error {
	result := s.db.WithContext(ctx).Where(
		"user_id = ? and source_id = ?", userID, sourceID,
	).Updates(newSubscription)
	if result.Error != nil {
		return result.Error
	}
	log.Debugf(
		"update %d row, userID %d sourceID %d new %#v", result.RowsAffected, userID, sourceID, newSubscription,
	)
	return nil
}

func (s *SubscriptionStorageImpl) UpsertSubscription(
	ctx context.Context, userID int64, sourceID uint, newSubscription *model.Subscribe,
) error {
	result := s.db.WithContext(ctx).Where(
		"user_id = ? and source_id = ?", userID, sourceID,
	).Save(newSubscription)
	if result.Error != nil {
		return result.Error
	}
	log.Debugf(
		"update %d row, userID %d sourceID %d new %#v", result.RowsAffected, userID, sourceID, newSubscription,
	)
	return nil
}
