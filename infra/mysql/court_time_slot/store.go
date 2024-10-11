package court_time_slot

import (
	"context"
	"fmt"
	"liga-app/entities"
	"time"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db}
}

func (s *Store) GetByCourtSlug(
	ctx context.Context,
	slug string,
	rangeStart time.Time,
	rangeEnd time.Time,
) ([]entities.CourtTimeSlot, error) {
	timeslots := []entities.CourtTimeSlot{}

	err := s.db.WithContext(ctx).
		Where("start >= ? AND end <= ?", rangeStart, rangeEnd).
		Find(&timeslots).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(ctx).Where(\"start >= ? AND end <= ?\", rangeStart, rangeEnd).Error : %w", err)
	}

	return timeslots, nil
}
