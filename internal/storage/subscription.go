package storage

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/indes/flowerss-bot/internal/model"
)

type SubscriptionStorageImpl struct {
	db *gorm.DB
}

func NewSubscriptionStorageImpl(db *gorm.DB) *SubscriptionStorageImpl {
	return &SubscriptionStorageImpl{db: db}
}

func (s *SubscriptionStorageImpl) AddSubscription(ctx context.Context, subscription *model.Subscribe) error {
	result := s.db.Create(subscription)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SubscriptionStorageImpl) GetSubscriptionsByUserID(
	ctx context.Context, userID int64, opts *GetSubscriptionsOptions,
) (*GetSubscriptionsResult, error) {
	var subscriptions []*model.Subscribe

	count := s.getSubscriptionsCount(opts)
	orderBy := s.getSubscriptionsOrderBy(opts)
	dbResult := s.db.Where(&model.Subscribe{UserID: userID}).Limit(count).Order(orderBy).Offset(opts.Offset).Find(subscriptions)
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, dbResult.Error
	}

	result := &GetSubscriptionsResult{}
	if len(subscriptions) > opts.Count {
		result.HasMore = true
		subscriptions = subscriptions[:opts.Count]
	}

	result.Subscriptions = subscriptions
	return result, nil
}

func (s *SubscriptionStorageImpl) GetSubscriptionsBySourceID(
	ctx context.Context, sourceID uint, opts *GetSubscriptionsOptions,
) (*GetSubscriptionsResult, error) {
	//var subscriptions []*model.Subscribe
	//s.db.Where("source_id=?", s.ID).Find(&subs)

	return nil, nil
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
		return "create_at desc"
	}
	return ""
}
