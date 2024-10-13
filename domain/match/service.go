package match

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"
	iService "liga-app/interfaces/services"
	"liga-app/utils"
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

func (s *Service) Get(ctx context.Context, p *common.Pagination, f *iService.MatchFilter) ([]entities.Match, int64, error) {
	repoFilter := &iRepo.MatchFilter{
		UserUUID: f.UserUUID,
	}

	startOfWeek := s.getStartOfWeek()

	if f != nil {
		if f.Status != nil {
			switch *f.Status {
			case iService.MATCH_FILTER_STATUS_UPCOMING:
				repoFilter.RoundIsAfter = &startOfWeek

			case iService.MATCH_FILTER_STATUS_COMPLETED:
				status := enums.MATCH_STATUS_COMPLETED
				repoFilter.Status = &status

			case iService.MATCH_FILTER_STATUS_OVERDUE:
				status := enums.MATCH_STATUS_WAITING
				repoFilter.Status = &status
				repoFilter.RoundIsBefore = &startOfWeek

			}
		}
	}

	matches, count, err := s.matchRepo.Get(ctx, p, repoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("s.matchRepo.Get(ctx, p, f) : %w", err)
	}

	for i, m := range matches {
		if m.Round.End != nil && startOfWeek.Before(*m.Round.End) && m.Status == string(enums.MATCH_STATUS_WAITING) {
			m.FilterStatus = string(iService.MATCH_FILTER_STATUS_UPCOMING)
			matches[i] = m
		} else if m.Round.End != nil && startOfWeek.After(*m.Round.End) && m.Status == string(enums.MATCH_STATUS_WAITING) {
			m.FilterStatus = string(iService.MATCH_FILTER_STATUS_OVERDUE)
			matches[i] = m
		} else {
			m.FilterStatus = string(m.Status)
			matches[i] = m
		}
	}

	return matches, count, nil
}

func (s *Service) getStartOfWeek() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	thisYear := utils.GetCurrentYear(loc)
	thisWeek := utils.GetCurrentWeek(loc)
	monday := utils.FirstDayOfISOWeek(thisYear, thisWeek, loc)
	return monday
}

func (s *Service) ProposeMatchSchedule(userUUID uuid.UUID, start time.Time, end time.Time) error {
	return nil
}

func (s *Service) AcceptMatchSchedule(matchScheduleRequestUUID uuid.UUID) error {
	return nil
}
