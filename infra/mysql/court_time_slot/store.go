package court_time_slot

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
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
		Joins("JOIN courts on courts.uuid = court_time_slots.court_uuid").
		Where("start >= ? AND end <= ?", rangeStart, rangeEnd).
		Where("courts.slug = ?", slug).
		Order("start ASC").
		Find(&timeslots).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(ctx).Where(\"start >= ? AND end <= ?\", rangeStart, rangeEnd).Error : %w", err)
	}

	return timeslots, nil
}

func (s *Store) Get(ctx context.Context, p *common.Pagination, f *repositories.CourtTimeSlotFilter) ([]entities.CourtTimeSlot, int64, error) {
	timeslots := []entities.CourtTimeSlot{}
	var count int64 = 0

	dbHandle := s.db.WithContext(ctx).Session(&gorm.Session{}).
		Preload("Court").
		Preload("Match").
		Offset(p.Offset).
		Limit(p.Limit)

	if f != nil {
		if len(f.CourtTimeSlotUUIDs) > 0 {
			dbHandle.Where("uuid in ?", f.CourtTimeSlotUUIDs)
		}
	}

	err := dbHandle.Find(&timeslots).Error
	if err != nil {
		return nil, 0, fmt.Errorf("dbHandle.Find(&timeslots).Error : %w", err)
	}

	dbCount := s.db.Model(&entities.CourtTimeSlot{}).
		WithContext(ctx).Session(&gorm.Session{})
	if f != nil {
		if len(f.CourtTimeSlotUUIDs) > 0 {
			dbHandle.Where("uuid in ?", f.CourtTimeSlotUUIDs)
		}
	}

	err = dbCount.Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("dbCount.Count(&count).Error : %w", err)
	}

	return timeslots, count, nil
}
