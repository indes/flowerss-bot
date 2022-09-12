package storage

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/indes/flowerss-bot/internal/model"
)

type SourceStorageImpl struct {
	db *gorm.DB
}

func NewSourceStorageImpl(db *gorm.DB) *SourceStorageImpl {
	return &SourceStorageImpl{db: db}
}

func (s *SourceStorageImpl) AddSource(ctx context.Context, source *model.Source) error {
	result := s.db.Create(source)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SourceStorageImpl) GetSource(ctx context.Context, id uint) (*model.Source, error) {
	var source = &model.Source{}
	result := s.db.Where(&model.Source{ID: id}).First(source)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, result.Error
	}

	return source, nil
}

func (s *SourceStorageImpl) GetSourceByURL(ctx context.Context, url string) (*model.Source, error) {
	var source = &model.Source{}
	result := s.db.Where(&model.Source{Link: url}).First(source)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, result.Error
	}
	return source, nil
}
