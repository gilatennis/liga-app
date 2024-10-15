package venue

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

func (s *Store) Get(c context.Context, p *common.Pagination, f *repositories.VenueFilter) ([]entities.Venue, int64, error) {
	venues := []entities.Venue{}
	var count int64

	dbHandle := s.db.WithContext(c).Model(&venues)

	if f != nil {
		// TODO:
		_ = 1
	}

	err := dbHandle.Find(&venues).
		Preload("Courts").
		Offset(p.Offset).
		Limit(p.Limit).
		Error
	if err != nil {
		return nil, 0, fmt.Errorf("dbHandle.Find(&venues)... : %w", err)
	}

	err = dbHandle.Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("dbHandle.Count(&count).Error : %w", err)
	}

	return venues, count, nil
}

func (s *Store) GetByUUID(c context.Context, venueUUID uuid.UUID) (*entities.Venue, error) {
	venue := entities.Venue{}

	err := s.db.WithContext(c).
		Preload("Courts").
		First(&venue, "uuid = ?", venueUUID).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(c).First(&venue...): %w", err)
	}

	return &venue, nil
}
