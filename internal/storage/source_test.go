package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indes/flowerss-bot/internal/model"
)

func TestSourceStorageImpl(t *testing.T) {
	db := GetTestDB(t)
	s := NewSourceStorageImpl(db)
	ctx := context.Background()
	s.Init(ctx)

	source := &model.Source{
		Link: "http://google.com",
	}

	t.Run(
		"add source", func(t *testing.T) {
			err := s.AddSource(ctx, source)
			assert.Nil(t, err)

			got, err := s.GetSource(ctx, source.ID)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, source.Link, got.Link)

			got, err = s.GetSourceByURL(ctx, source.Link)
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, source.ID, got.ID)

			err = s.Delete(ctx, got.ID)
			assert.Nil(t, err)

			got, err = s.GetSource(ctx, source.ID)
			assert.Equal(t, ErrRecordNotFound, err)
			assert.Nil(t, got)
		},
	)

}
