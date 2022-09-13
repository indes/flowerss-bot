package core

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indes/flowerss-bot/internal/storage/mock"
)

type mockStorage struct {
	User         *mock.MockUser
	Content      *mock.MockContent
	Source       *mock.MockSource
	Subscription *mock.MockSubscription
	Ctrl         *gomock.Controller
}

func getTestCore(t *testing.T) (*Core, *mockStorage) {
	ctrl := gomock.NewController(t)

	s := &mockStorage{
		Subscription: mock.NewMockSubscription(ctrl),
		User:         mock.NewMockUser(ctrl),
		Content:      mock.NewMockContent(ctrl),
		Source:       mock.NewMockSource(ctrl),
		Ctrl:         ctrl,
	}
	c := NewCore(s.User, s.Content, s.Source, s.Subscription)
	return c, s
}

func TestCore_AddSubscription(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()

	userID := int64(1)
	sourceID := uint(101)
	t.Run(
		"exist error", func(t *testing.T) {
			s.Subscription.EXPECT().SubscriptionExist(ctx, userID, sourceID).Return(false, errors.New("err")).Times(1)
			err := c.AddSubscription(ctx, userID, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"exist subscription", func(t *testing.T) {
			s.Subscription.EXPECT().SubscriptionExist(ctx, userID, sourceID).Return(true, nil).Times(1)
			err := c.AddSubscription(ctx, userID, sourceID)
			assert.Equal(t, ErrSubscriptionExist, err)
		},
	)

	t.Run(
		"subscribe fail", func(t *testing.T) {
			s.Subscription.EXPECT().SubscriptionExist(ctx, userID, sourceID).Return(false, nil).Times(1)
			s.Subscription.EXPECT().AddSubscription(ctx, gomock.Any()).Return(errors.New("err")).Times(1)

			err := c.AddSubscription(ctx, userID, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"subscribe ok", func(t *testing.T) {
			s.Subscription.EXPECT().SubscriptionExist(ctx, userID, sourceID).Return(false, nil).Times(1)
			s.Subscription.EXPECT().AddSubscription(ctx, gomock.Any()).Return(nil).Times(1)

			err := c.AddSubscription(ctx, userID, sourceID)
			assert.Nil(t, err)
		},
	)
}
