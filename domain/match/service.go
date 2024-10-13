package match

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	matchRepo         iRepo.IMatchRepository
	courtTimeSlotRepo iRepo.ICourtTimeSlotRepository
}

func NewService(
	matchRepo iRepo.IMatchRepository,
	courtTimeSlotRepo iRepo.ICourtTimeSlotRepository,
) *Service {
	return &Service{
		matchRepo:         matchRepo,
		courtTimeSlotRepo: courtTimeSlotRepo,
	}
}

func (s *Service) Get(ctx context.Context, p *common.Pagination, f *iRepo.MatchFilter) ([]entities.Match, int64, error) {
	matches, count, err := s.matchRepo.Get(ctx, p, f)
	if err != nil {
		return nil, 0, fmt.Errorf("s.matchRepo.Get(ctx, p, f) : %w", err)
	}

	return matches, count, nil
}

func (s *Service) ProposeMatchSchedule(userUUID uuid.UUID, start time.Time, end time.Time) error {
	return nil
}

func (s *Service) AcceptMatchSchedule(matchScheduleRequestUUID uuid.UUID) error {
	return nil
}
