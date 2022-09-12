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
	return &ContentStorageImpl{db: db.Model(&model.Content{})}
}

func (s *ContentStorageImpl) Init(ctx context.Context) error {
	return s.db.Migrator().AutoMigrate(&model.Content{})
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

func (s *ContentStorageImpl) HashIDExist(ctx context.Context, hashID string) (bool, error) {
	var count int64
	result := s.db.WithContext(ctx).Where("hash_id = ?", hashID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return (count > 0), nil
}
