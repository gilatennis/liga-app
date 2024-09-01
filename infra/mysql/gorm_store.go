package mysql

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db}
}

func (s *GormStore) Save(ctx context.Context, entity interfaces.IEntity, user *entities.User) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Save(entity).
		Error
	if err != nil {
		return nil, err
	}

	return entity.GetUUID(), nil
}

func (s *GormStore) Delete(ctx context.Context, entity interfaces.IEntity, user *entities.User) (*uuid.UUID, error) {
	return nil, nil
}
