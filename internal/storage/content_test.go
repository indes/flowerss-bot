package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indes/flowerss-bot/internal/model"
)

func TestContentStorageImpl(t *testing.T) {
	db := GetTestDB(t)
	s := NewContentStorageImpl(db)
	ctx := context.Background()
	s.Init(ctx)

	content := &model.Content{
		SourceID: 1,
		HashID:   "id",
	}
	content2 := &model.Content{
		SourceID: 1,
		HashID:   "id2",
	}

	t.Run(
		"add content", func(t *testing.T) {
			err := s.AddContent(ctx, content)
			assert.Nil(t, err)
			err = s.AddContent(ctx, content2)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"hash id exist", func(t *testing.T) {
			exist, err := s.HashIDExist(ctx, content.HashID)
			assert.Nil(t, err)
			assert.True(t, exist)
		},
	)

	t.Run(
		"del content", func(t *testing.T) {
			got, err := s.DeleteSourceContents(ctx, content.SourceID)
			assert.Nil(t, err)
			assert.Equal(t, int64(2), got)
		},
	)

	t.Run(
		"hash id exist2", func(t *testing.T) {
			exist, err := s.HashIDExist(ctx, content.HashID)
			assert.Nil(t, err)
			assert.False(t, exist)
		},
	)
}
