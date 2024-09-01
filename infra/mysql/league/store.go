package league

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db}
}

func (s *Store) Get(ctx context.Context, p *common.Pagination, f *repositories.LeagueFilter) ([]entities.League, int64, error) {
	leagues := []entities.League{}
	var count int64
	err := s.db.WithContext(ctx).
		Preload("Wilayah").
		Find(&leagues).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}

	return leagues, count, nil
}

func (s *Store) GetByUUID(ctx context.Context, leagueUUID uuid.UUID) (*entities.League, error) {
	return nil, nil
}

func (s *Store) Create(ctx context.Context, e *entities.League) (*uuid.UUID, error) {
	return nil, nil
}

func (s *Store) Update(ctx context.Context, e *entities.League) (*uuid.UUID, error) {
	return nil, nil
}

func (s *Store) Delete(ctx context.Context, e *entities.League) (*uuid.UUID, error) {
	return nil, nil
}
