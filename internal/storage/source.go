package storage

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/indes/flowerss-bot/internal/model"
)

type SourceStorageImpl struct {
	db *gorm.DB
}

func NewSourceStorageImpl(db *gorm.DB) *SourceStorageImpl {
	return &SourceStorageImpl{db: db}
}

func (s SourceStorageImpl) AddSource(ctx context.Context, source *model.Source) error {
	result := s.db.Create(source)
	if len(result.GetErrors()) != 0 {
		return result.GetErrors()[0]
	}
	return nil
}

func (s SourceStorageImpl) GetSource(ctx context.Context, id uint) (*model.Source, error) {
	var Source = &model.Source{}
	result := s.db.Where(&model.Source{ID: id}).First(Source)
	if result.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	if len(result.GetErrors()) != 0 {
		return nil, result.GetErrors()[0]
	}
	return Source, nil
}

func (s SourceStorageImpl) GetSourceByURL(ctx context.Context, url string) (*model.Source, error) {
	var Source = &model.Source{}
	result := s.db.Where(&model.Source{Link: url}).First(Source)
	if result.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	if len(result.GetErrors()) != 0 {
		return nil, result.GetErrors()[0]
	}
	return Source, nil
}
