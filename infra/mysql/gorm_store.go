package mysql

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces"
	iCommon "liga-app/interfaces/common"

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

func (s *GormStore) Get(
	ctx context.Context,
	p *iCommon.Pagination,
	models any,
	where *map[string]interface{},
) (any, int64, error) {
	var count int64

	dbHandle := s.db.WithContext(ctx).
		Limit(p.Limit).
		Offset(p.Offset)

	if where != nil {
		dbHandle = dbHandle.Where(*where)
	}

	err := dbHandle.
		Find(&models).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}

	return models, count, nil
}
