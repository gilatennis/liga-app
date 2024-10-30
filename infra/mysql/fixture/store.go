package fixture

import (
	"context"
	"fmt"
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

func (s *Store) Get(ctx context.Context, p *common.Pagination, f *repositories.FixtureFilter) ([]entities.Fixture, int64, error) {
	fixtures := []entities.Fixture{}
	var count int64

	dbHandle := s.db.WithContext(ctx).
		Offset(p.Offset).
		Limit(p.Limit)

	dbCount := s.db.WithContext(ctx).
		Model(&entities.Fixture{})

	if f != nil {
		if f.LeagueUUID != nil {
			dbHandle = dbHandle.Where("league_uuid = ?", *f.LeagueUUID)
			dbCount = dbCount.Where("league_uuid = ?", *f.LeagueUUID)
		}
	}

	err := dbHandle.Find(&fixtures).Error
	if err != nil {
		return nil, 0, fmt.Errorf("dbHandle.Find(&fixtures).Error : %w", err)
	}

	err = dbCount.Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("dbCount.Count(&count).Error : %w", err)
	}

	return fixtures, count, nil
}

func (s *Store) GetByUUID(ctx context.Context, fixtureUUID uuid.UUID) (*entities.Fixture, error) {
	f := entities.Fixture{}
	err := s.db.WithContext(ctx).
		Preload("Groups").
		First(&f, "uuid = ?", fixtureUUID).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db...First(&f, \"uuid = ?\", fixtureUUID=%s).Error : %w", fixtureUUID, err)
	}

	return &f, nil
}
