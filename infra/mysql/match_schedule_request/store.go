package match_schedule_request

import (
	"context"
	"errors"
	"fmt"
	"liga-app/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db}
}

func (s *Store) Create(ctx context.Context, req *entities.MatchScheduleRequest) (*uuid.UUID, error) {
	tx := s.db.WithContext(ctx).Begin()

	existingRequest := entities.MatchScheduleRequest{}
	err := tx.Where("match_uuid = ?", req.MatchUUID).First(&existingRequest).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("tx.Where(\"match_uuid = ?\", req.MatchUUID).First(&existingRequest).Error")
	}

	if existingRequest.UUID != uuid.Nil {
		err = tx.Unscoped().
			Delete(&entities.MatchScheduleRequestCourtTimeSlot{}, "match_schedule_request_uuid = ?", existingRequest.UUID).
			Error
		if err != nil {
			return nil, fmt.Errorf("tx.Unscoped().Delete(&entities.MatchScheduleRequestCourtTimeSlot{}).Error : %w", err)
		}

		err = tx.Unscoped().
			Delete(&entities.MatchScheduleRequest{}, "uuid = ?", existingRequest.UUID).
			Error
		if err != nil {
			return nil, fmt.Errorf("tx.Unscoped().Delete(&entities.MatchScheduleRequest{}).Error : %w", err)
		}
	}

	err = tx.Omit(clause.Associations).Create(req).Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(ctx).Create(req).Error")
	}

	reqTimeSlots := []entities.MatchScheduleRequestCourtTimeSlot{}
	for _, ts := range req.TimeSlots {
		rt := entities.MatchScheduleRequestCourtTimeSlot{
			MatchScheduleRequestUUID: req.UUID,
			CourtTimeSlotUUID:        ts.UUID,
		}
		reqTimeSlots = append(reqTimeSlots, rt)
	}

	err = tx.Omit(clause.Associations).Create(reqTimeSlots).Error
	if err != nil {
		return nil, fmt.Errorf("tx.Omit(clause.Associations).Create(reqTimeSlots).Error")
	}

	err = tx.Model(&entities.Match{}).
		Where("uuid = ?", req.MatchUUID).
		Update("match_schedule_request_uuid", req.UUID).
		Error
	if err != nil {
		return nil, fmt.Errorf("tx...Update(match_schedule_request_uuid, ..): %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, fmt.Errorf("tx.Commit().Error : %w", err)
	}

	return &req.UUID, err
}

func (s *Store) GetByUUID(ctx context.Context, reqUUID uuid.UUID) (*entities.MatchScheduleRequest, error) {
	req := entities.MatchScheduleRequest{}

	err := s.db.WithContext(ctx).
		Where("uuid = ?", reqUUID).
		Preload("Match").
		Preload("Match.TeamA").
		Preload("Match.TeamA.Captain").
		Preload("Match.TeamB").
		Preload("Match.TeamB.Captain").
		Preload("Match.MatchScheduleRequest").
		First(&req).Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(ctx).Where(\"uuid = ?\", reqUUID).First(&req).Error : %w", err)
	}

	return &req, nil
}

func (s *Store) Update(ctx context.Context, req *entities.MatchScheduleRequest) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Omit(clause.Associations).
		Save(req).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db...Save(req).Error : %w", err)
	}

	return &req.UUID, nil
}
