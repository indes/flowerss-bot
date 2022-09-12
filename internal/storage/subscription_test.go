package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indes/flowerss-bot/internal/model"
)

func TestSubscriptionStorageImpl(t *testing.T) {
	db := GetTestDB(t)
	s := NewSubscriptionStorageImpl(db)
	ctx := context.Background()
	s.Init(ctx)

	subscriptions := []*model.Subscribe{
		&model.Subscribe{
			SourceID:           1,
			UserID:             100,
			EnableNotification: 1,
		},
		&model.Subscribe{
			SourceID:           1,
			UserID:             101,
			EnableNotification: 1,
		},

		&model.Subscribe{
			SourceID:           2,
			UserID:             100,
			EnableNotification: 1,
		},
		&model.Subscribe{
			SourceID:           2,
			UserID:             101,
			EnableNotification: 1,
		},
		&model.Subscribe{
			SourceID:           3,
			UserID:             101,
			EnableNotification: 1,
		},
	}

	t.Run(
		"add subscription", func(t *testing.T) {
			for _, subscription := range subscriptions {
				err := s.AddSubscription(ctx, subscription)
				assert.Nil(t, err)
			}
			got, err := s.CountSubscriptions(ctx)
			assert.Nil(t, err)
			assert.Equal(t, int64(5), got)

			exist, err := s.SubscriptionExist(ctx, 101, 1)
			assert.Nil(t, err)
			assert.True(t, exist)

			opt := &GetSubscriptionsOptions{
				Count: 2,
			}
			result, err := s.GetSubscriptionsByUserID(ctx, 101, opt)
			assert.Nil(t, err)
			assert.Equal(t, 2, len(result.Subscriptions))
			assert.True(t, result.HasMore)

			opt = &GetSubscriptionsOptions{
				Count:  1,
				Offset: 2,
			}
			result, err = s.GetSubscriptionsByUserID(ctx, 101, opt)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(result.Subscriptions))
			assert.False(t, result.HasMore)

			opt = &GetSubscriptionsOptions{
				Count: 2,
			}
			result, err = s.GetSubscriptionsBySourceID(ctx, 1, opt)
			assert.Nil(t, err)
			assert.Equal(t, 2, len(result.Subscriptions))
			assert.False(t, result.HasMore)

			opt = &GetSubscriptionsOptions{
				Count:  1,
				Offset: 2,
			}
			result, err = s.GetSubscriptionsByUserID(ctx, 1, opt)
			assert.Nil(t, err)
			assert.Equal(t, 0, len(result.Subscriptions))
			assert.False(t, result.HasMore)

			got, err = s.DeleteSubscription(ctx, 101, 1)
			assert.Nil(t, err)
			assert.Equal(t, int64(1), got)

			exist, err = s.SubscriptionExist(ctx, 101, 1)
			assert.Nil(t, err)
			assert.False(t, exist)

			got, err = s.CountSubscriptions(ctx)
			assert.Nil(t, err)
			assert.Equal(t, int64(4), got)
		},
	)
}
