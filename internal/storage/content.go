package storage

import (
	"context"

	"gorm.io/gorm"

	"github.com/indes/flowerss-bot/internal/model"
)

type ContentStorageImpl struct {
	db *gorm.DB
}

func NewContentStorageImpl(db *gorm.DB) *ContentStorageImpl {
	return &ContentStorageImpl{db: db}
}

func (s *ContentStorageImpl) DeleteSourceContents(ctx context.Context, sourceID uint) (int64, error) {
	result := s.db.WithContext(ctx).Where("source_id = ?", sourceID).Delete(&model.Content{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (s *ContentStorageImpl) AddContent(ctx context.Context, content *model.Content) error {
	result := s.db.WithContext(ctx).Create(content)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
