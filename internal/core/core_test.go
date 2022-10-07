package core

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/indes/flowerss-bot/internal/model"
	"github.com/indes/flowerss-bot/internal/storage"
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
	c := NewCore(s.User, s.Content, s.Source, s.Subscription, nil, nil)
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

func TestCore_GetUserSubscribedSources(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()

	userID := int64(1)
	sourceID1 := uint(101)
	sourceID2 := uint(102)
	subscriptionsResult := &storage.GetSubscriptionsResult{
		Subscriptions: []*model.Subscribe{
			&model.Subscribe{SourceID: sourceID1},
			&model.Subscribe{SourceID: sourceID2},
		},
	}
	t.Run(
		"subscription err", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscriptionsByUserID(ctx, userID, gomock.Any()).Return(
				nil, errors.New("err"),
			)

			sources, err := c.GetUserSubscribedSources(ctx, userID)
			assert.Error(t, err)
			assert.Nil(t, sources)
		},
	)

	t.Run(
		"source err", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscriptionsByUserID(ctx, userID, gomock.Any()).Return(
				subscriptionsResult, nil,
			)

			s.Source.EXPECT().GetSource(ctx, sourceID1).Return(
				nil, errors.New("err"),
			).Times(1)
			s.Source.EXPECT().GetSource(ctx, gomock.Any()).Return(
				&model.Source{}, nil,
			)

			sources, err := c.GetUserSubscribedSources(ctx, userID)
			assert.Nil(t, err)
			assert.Equal(t, len(subscriptionsResult.Subscriptions)-1, len(sources))
		},
	)

	t.Run(
		"source success", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscriptionsByUserID(ctx, userID, gomock.Any()).Return(
				subscriptionsResult, nil,
			)

			s.Source.EXPECT().GetSource(ctx, gomock.Any()).Return(
				&model.Source{}, nil,
			).Times(len(subscriptionsResult.Subscriptions))

			sources, err := c.GetUserSubscribedSources(ctx, userID)
			assert.Nil(t, err)
			assert.Equal(t, len(subscriptionsResult.Subscriptions), len(sources))
		},
	)
}

func TestCore_Unsubscribe(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()

	userID := int64(1)
	sourceID1 := uint(101)

	t.Run(
		"SubscriptionExist err", func(t *testing.T) {
			s.Subscription.EXPECT().SubscriptionExist(ctx, userID, sourceID1).Return(
				false, errors.New("err"),
			).Times(1)
			err := c.Unsubscribe(ctx, userID, sourceID1)
			assert.Error(t, err)
		},
	)

	t.Run(
		"subscription not exist", func(t *testing.T) {
			s.Subscription.EXPECT().SubscriptionExist(ctx, userID, sourceID1).Return(
				false, nil,
			).Times(1)
			err := c.Unsubscribe(ctx, userID, sourceID1)
			assert.Equal(t, ErrSubscriptionNotExist, err)
		},
	)

	s.Subscription.EXPECT().SubscriptionExist(ctx, gomock.Any(), gomock.Any()).Return(
		true, nil,
	).AnyTimes()

	t.Run(
		"unsubscribe failed", func(t *testing.T) {
			s.Subscription.EXPECT().DeleteSubscription(ctx, userID, sourceID1).Return(
				int64(1), errors.New("err"),
			).Times(1)
			err := c.Unsubscribe(ctx, userID, sourceID1)
			assert.Error(t, err)
		},
	)

	s.Subscription.EXPECT().DeleteSubscription(ctx, gomock.Any(), gomock.Any()).Return(
		int64(1), nil,
	).AnyTimes()

	t.Run(
		"count subs", func(t *testing.T) {
			s.Subscription.EXPECT().CountSourceSubscriptions(ctx, sourceID1).Return(
				int64(1), errors.New("err"),
			).Times(1)
			err := c.Unsubscribe(ctx, userID, sourceID1)
			assert.Error(t, err)

			s.Subscription.EXPECT().CountSourceSubscriptions(ctx, sourceID1).Return(
				int64(1), nil,
			).Times(1)
			err = c.Unsubscribe(ctx, userID, sourceID1)
			assert.Nil(t, err)
		},
	)

	s.Subscription.EXPECT().CountSourceSubscriptions(ctx, gomock.Any()).Return(
		int64(0), nil,
	).AnyTimes()

	t.Run(
		"remove source", func(t *testing.T) {
			s.Source.EXPECT().Delete(ctx, sourceID1).Return(
				errors.New("err"),
			).Times(1)

			err := c.Unsubscribe(ctx, userID, sourceID1)
			assert.Error(t, err)

			s.Source.EXPECT().Delete(ctx, sourceID1).Return(nil).AnyTimes()

			s.Content.EXPECT().DeleteSourceContents(ctx, sourceID1).Return(int64(0), errors.New("err")).Times(1)
			err = c.Unsubscribe(ctx, userID, sourceID1)
			assert.Error(t, err)

			s.Content.EXPECT().DeleteSourceContents(ctx, sourceID1).Return(int64(1), nil).Times(1)
			err = c.Unsubscribe(ctx, userID, sourceID1)
			assert.Nil(t, err)
		},
	)
}

func TestCore_GetSourceByURL(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()
	sourceURL := "http://google.com"

	t.Run(
		"source err", func(t *testing.T) {
			s.Source.EXPECT().GetSourceByURL(ctx, sourceURL).Return(
				nil, errors.New("err"),
			).Times(1)
			got, err := c.GetSourceByURL(ctx, sourceURL)
			assert.Error(t, err)
			assert.Nil(t, got)
		},
	)

	t.Run(
		"source not exist", func(t *testing.T) {
			s.Source.EXPECT().GetSourceByURL(ctx, sourceURL).Return(
				nil, storage.ErrRecordNotFound,
			).Times(1)
			got, err := c.GetSourceByURL(ctx, sourceURL)
			assert.Equal(t, ErrSourceNotExist, err)
			assert.Nil(t, got)
		},
	)

	t.Run(
		"ok", func(t *testing.T) {
			s.Source.EXPECT().GetSourceByURL(ctx, sourceURL).Return(
				&model.Source{}, nil,
			).Times(1)
			got, err := c.GetSourceByURL(ctx, sourceURL)
			assert.Nil(t, err)
			assert.NotNil(t, got)
		},
	)
}

func TestCore_GetSource(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()
	sourceID := uint(1)

	t.Run(
		"source err", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				nil, errors.New("err"),
			).Times(1)
			got, err := c.GetSource(ctx, sourceID)
			assert.Error(t, err)
			assert.Nil(t, got)
		},
	)

	t.Run(
		"source not exist", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				nil, storage.ErrRecordNotFound,
			).Times(1)
			got, err := c.GetSource(ctx, sourceID)
			assert.Equal(t, ErrSourceNotExist, err)
			assert.Nil(t, got)
		},
	)

	t.Run(
		"ok", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				&model.Source{}, nil,
			).Times(1)
			got, err := c.GetSource(ctx, sourceID)
			assert.Nil(t, err)
			assert.NotNil(t, got)
		},
	)
}

func TestCore_GetSubscription(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()
	userID := int64(101)
	sourceID := uint(1)

	t.Run(
		"subscription err", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscription(ctx, userID, sourceID).Return(
				nil, errors.New("err"),
			).Times(1)
			got, err := c.GetSubscription(ctx, userID, sourceID)
			assert.Error(t, err)
			assert.Nil(t, got)
		},
	)

	t.Run(
		"subscription not exist", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscription(ctx, userID, sourceID).Return(
				nil, storage.ErrRecordNotFound,
			).Times(1)
			got, err := c.GetSubscription(ctx, userID, sourceID)
			assert.Equal(t, ErrSubscriptionNotExist, err)
			assert.Nil(t, got)
		},
	)

	t.Run(
		"ok", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscription(ctx, userID, sourceID).Return(
				&model.Subscribe{}, nil,
			).Times(1)
			got, err := c.GetSubscription(ctx, userID, sourceID)
			assert.Nil(t, err)
			assert.NotNil(t, got)
		},
	)
}

func TestCore_DisableSourceUpdate(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()
	sourceID := uint(1)

	t.Run(
		"get source err", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				nil, errors.New("err"),
			).Times(1)
			err := c.DisableSourceUpdate(ctx, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"update source err", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				&model.Source{}, nil,
			).Times(1)

			s.Source.EXPECT().UpsertSource(ctx, sourceID, gomock.Any()).Return(
				errors.New("err"),
			).Times(1)
			err := c.DisableSourceUpdate(ctx, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"update source err", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				&model.Source{}, nil,
			).Times(1)

			s.Source.EXPECT().UpsertSource(ctx, sourceID, gomock.Any()).Return(
				nil,
			).Times(1)
			err := c.DisableSourceUpdate(ctx, sourceID)
			assert.Nil(t, err)
		},
	)
}

func TestCore_ClearSourceErrorCount(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()
	sourceID := uint(1)

	t.Run(
		"get source err", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				nil, errors.New("err"),
			).Times(1)
			err := c.ClearSourceErrorCount(ctx, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"update source err", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				&model.Source{}, nil,
			).Times(1)

			s.Source.EXPECT().UpsertSource(ctx, sourceID, gomock.Any()).Return(
				errors.New("err"),
			).Times(1)
			err := c.ClearSourceErrorCount(ctx, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"update source err", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				&model.Source{}, nil,
			).Times(1)

			s.Source.EXPECT().UpsertSource(ctx, sourceID, gomock.Any()).Return(
				nil,
			).Times(1)
			err := c.ClearSourceErrorCount(ctx, sourceID)
			assert.Nil(t, err)
		},
	)
}

func TestCore_ToggleSubscriptionNotice(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()
	userID := int64(123)
	sourceID := uint(1)

	t.Run(
		"get subscription err", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscription(ctx, userID, sourceID).Return(
				nil, errors.New("err"),
			).Times(1)
			err := c.ToggleSubscriptionNotice(ctx, userID, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"update subscription err", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscription(ctx, userID, sourceID).Return(
				&model.Subscribe{}, nil,
			).Times(1)

			s.Subscription.EXPECT().UpsertSubscription(ctx, userID, sourceID, gomock.Any()).Return(
				errors.New("err"),
			).Times(1)

			err := c.ToggleSubscriptionNotice(ctx, userID, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"ok", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscription(ctx, userID, sourceID).Return(
				&model.Subscribe{}, nil,
			).Times(1)

			s.Subscription.EXPECT().UpsertSubscription(ctx, userID, sourceID, gomock.Any()).Return(
				nil,
			).Times(1)

			err := c.ToggleSubscriptionNotice(ctx, userID, sourceID)
			assert.Nil(t, err)
		},
	)
}

func TestCore_ToggleSourceUpdateStatus(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()
	sourceID := uint(1)

	t.Run(
		"get source err", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				nil, errors.New("err"),
			).Times(1)
			err := c.ToggleSourceUpdateStatus(ctx, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"update source err", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				&model.Source{}, nil,
			).Times(1)

			s.Source.EXPECT().UpsertSource(ctx, sourceID, gomock.Any()).Return(
				errors.New("err"),
			).Times(1)
			err := c.ToggleSourceUpdateStatus(ctx, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"ok", func(t *testing.T) {
			s.Source.EXPECT().GetSource(ctx, sourceID).Return(
				&model.Source{}, nil,
			).Times(1)

			s.Source.EXPECT().UpsertSource(ctx, sourceID, gomock.Any()).Return(
				nil,
			).Times(1)
			err := c.ToggleSourceUpdateStatus(ctx, sourceID)
			assert.Nil(t, err)
		},
	)
}

func TestCore_ToggleSubscriptionTelegraph(t *testing.T) {
	c, s := getTestCore(t)
	defer s.Ctrl.Finish()
	ctx := context.Background()
	userID := int64(123)
	sourceID := uint(1)

	t.Run(
		"get subscription err", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscription(ctx, userID, sourceID).Return(
				nil, errors.New("err"),
			).Times(1)
			err := c.ToggleSubscriptionTelegraph(ctx, userID, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"update subscription err", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscription(ctx, userID, sourceID).Return(
				&model.Subscribe{}, nil,
			).Times(1)

			s.Subscription.EXPECT().UpsertSubscription(ctx, userID, sourceID, gomock.Any()).Return(
				errors.New("err"),
			).Times(1)

			err := c.ToggleSubscriptionTelegraph(ctx, userID, sourceID)
			assert.Error(t, err)
		},
	)

	t.Run(
		"ok", func(t *testing.T) {
			s.Subscription.EXPECT().GetSubscription(ctx, userID, sourceID).Return(
				&model.Subscribe{}, nil,
			).Times(1)

			s.Subscription.EXPECT().UpsertSubscription(ctx, userID, sourceID, gomock.Any()).Return(
				nil,
			).Times(1)

			err := c.ToggleSubscriptionTelegraph(ctx, userID, sourceID)
			assert.Nil(t, err)
		},
	)
}
