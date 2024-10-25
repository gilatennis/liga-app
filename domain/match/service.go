package match

import (
	"context"
	"fmt"
	domainErrors "liga-app/domain/errors"
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
	matchRepo                iRepo.IMatchRepository
	courtTimeSlotRepo        iRepo.ICourtTimeSlotRepository
	matchScheduleRequestRepo iRepo.IMatchScheduleRequestRepository
	matchSetRepo             iRepo.IMatchSetRepository
}

func NewService(
	matchRepo iRepo.IMatchRepository,
	courtTimeSlotRepo iRepo.ICourtTimeSlotRepository,
	matchScheduleRequestRepo iRepo.IMatchScheduleRequestRepository,
	matchSetRepo iRepo.IMatchSetRepository,
) *Service {
	return &Service{
		matchRepo:                matchRepo,
		courtTimeSlotRepo:        courtTimeSlotRepo,
		matchScheduleRequestRepo: matchScheduleRequestRepo,
		matchSetRepo:             matchSetRepo,
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
		s.setMatchFilterStatus(&m)
		matches[i] = m
		/*
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
		*/
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

func (s *Service) setMatchFilterStatus(m *entities.Match) {
	if m.IsComplete && m.IsCompleteConfirmed {
		m.FilterStatus = string(enums.MATCH_STATUS_COMPLETED)
		return
	}

	startOfWeek := s.getStartOfWeek()
	if m.Round.End != nil && startOfWeek.Before(*m.Round.End) && m.Status == string(enums.MATCH_STATUS_WAITING) {
		m.FilterStatus = string(iService.MATCH_FILTER_STATUS_UPCOMING)
	} else if m.Round.End != nil && startOfWeek.After(*m.Round.End) && m.Status == string(enums.MATCH_STATUS_WAITING) {
		m.FilterStatus = string(iService.MATCH_FILTER_STATUS_OVERDUE)
	} else {
		m.FilterStatus = string(m.Status)
	}

	if m.TeamA == nil || m.TeamB == nil {
		m.FilterStatus = "NO SCHEDULE"
	}
}

func (s *Service) GetByUUID(ctx context.Context, matchUUID uuid.UUID) (*entities.Match, error) {
	match, err := s.matchRepo.GetByUUID(ctx, matchUUID)
	if err != nil {
		return nil, fmt.Errorf("s.matchRepo.GetByUUID(ctx, matchUUID=%s): %w", matchUUID, err)
	}

	s.setMatchFilterStatus(match)

	return match, nil
}

func (s *Service) RequestMatchSchedule(ctx context.Context, userUUID uuid.UUID, matchUUID uuid.UUID, timeslotUUIDs []uuid.UUID) (*uuid.UUID, error) {
	p := common.NewPagination(1, len(timeslotUUIDs), []common.OrderBy{})
	f := &iRepo.CourtTimeSlotFilter{CourtTimeSlotUUIDs: timeslotUUIDs}
	timeslots, _, err := s.courtTimeSlotRepo.Get(ctx, p, f)
	if err != nil {
		return nil, fmt.Errorf("s.courtTimeSlotRepo.Get(ctx, p, f) : %w", err)
	}

	if len(timeslots) == 0 {
		return nil, &domainErrors.ErrNotFound{Message: "timeslots not found"}
	}

	courtUUID := timeslots[0].CourtUUID

	for _, ts := range timeslots {
		if ts.CourtUUID != courtUUID {
			return nil, &domainErrors.ErrBadRequest{Message: "timeslots have to be for 1 court only"}
		}
		if ts.Match != nil {
			return nil, &domainErrors.ErrConflict{Message: "timeslots requested already have match(es) assigned"}
		}
	}

	req := entities.MatchScheduleRequest{
		Base: entities.Base{
			CreatedBy: userUUID,
			UpdatedBy: userUUID,
		},
		MatchUUID: matchUUID,
		CourtUUID: &courtUUID,
		TimeSlots: timeslots,
	}

	reqUUID, err := s.matchScheduleRequestRepo.Create(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("s.matchScheduleRequestRepo.Create(ctx, &req): %w", err)
	}

	return reqUUID, nil
}

func (s *Service) RequestMatchScheduleOwnVenue(ctx context.Context, userUUID uuid.UUID, matchUUID uuid.UUID, req *entities.MatchScheduleRequest) (*uuid.UUID, error) {
	req.CreatedBy = userUUID
	req.UpdatedBy = userUUID
	req.MatchUUID = matchUUID

	reqUUID, err := s.matchScheduleRequestRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("s.matchScheduleRequestRepo.Create(ctx, req): %w", err)
	}

	return reqUUID, nil
}

func (s *Service) ConfirmMatchSchedule(ctx context.Context, userUUID uuid.UUID, matchScheduleRequestUUID uuid.UUID) (*uuid.UUID, error) {
	req, err := s.matchScheduleRequestRepo.GetByUUID(ctx, matchScheduleRequestUUID)
	if err != nil {
		return nil, fmt.Errorf("s.matchScheduleRequestRepo.GetByUUID(ctx, matchScheduleRequestUUID=%s) : %w", matchScheduleRequestUUID, err)
	}

	if !isMatchScheduleConfirmer(userUUID, &req.Match) {
		return nil, &domainErrors.ErrUnauthorized{Message: "not authorized to confirm schedule"}
	}

	req.IsAccepted = true
	req.UpdatedBy = userUUID

	updatedUUID, err := s.matchScheduleRequestRepo.AcceptMatchScheduleRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("s.matchScheduleRequestRepo.Update(ctx, req) : %w", err)
	}

	return updatedUUID, nil
}

func isMatchScheduleConfirmer(userUUID uuid.UUID, match *entities.Match) bool {
	return isMatchParticipantCaptain(userUUID, match) && userUUID != match.MatchScheduleRequest.CreatedBy
}

func isMatchParticipantCaptain(userUUID uuid.UUID, match *entities.Match) bool {
	return userUUID == match.TeamA.CaptainUserUUID || userUUID == match.TeamB.CaptainUserUUID
}

func (s *Service) SaveMatchSet(ctx context.Context, userUUID uuid.UUID, matchSet *entities.MatchSet) (*uuid.UUID, error) {
	match, err := s.matchRepo.GetByUUID(ctx, matchSet.MatchUUID)
	if err != nil {
		return nil, fmt.Errorf("s.matchRepo.GetByUUID(ctx, matchSet.MatchUUID=%s) : %w", matchSet.MatchUUID, err)
	}

	if match.TeamA.CaptainUserUUID != userUUID && match.TeamB.CaptainUserUUID != userUUID {
		return nil, &domainErrors.ErrUnauthorized{Message: "not authorized to save match set"}
	}

	matchSetUUID, err := s.matchSetRepo.SaveMatchSet(ctx, userUUID, matchSet)
	if err != nil {
		return nil, fmt.Errorf("s.matchSetRepo.SaveMatchSet(ctx, matchSet) : %w", err)
	}

	return matchSetUUID, nil
}

func (s *Service) CompleteMatch(ctx context.Context, userUUID, matchUUID uuid.UUID) (*uuid.UUID, error) {
	match, err := s.matchRepo.GetByUUID(ctx, matchUUID)
	if err != nil {
		return nil, fmt.Errorf("s.matchRepo.GetByUUID(ctx, matchUUID=%s) : %w", matchUUID, err)
	}

	if match.TeamA.CaptainUserUUID != userUUID && match.TeamB.CaptainUserUUID != userUUID {
		return nil, &domainErrors.ErrUnauthorized{Message: "not authorized to complete match"}
	}

	if match.League.ParticipationCategory == enums.CATEGORY_PARTICIPATION_TEAMS {
		if len(match.Sets) < 3 {
			return nil, &domainErrors.ErrBadRequest{Message: "skor pertandingan belum lengkap"}
		}
		for _, s := range match.Sets {
			if s.TeamAScore == 0 && s.TeamBScore == 0 {
				return nil, &domainErrors.ErrBadRequest{Message: "skor pertandingan belum lengkap"}
			}
		}

	} else {
		if len(match.Sets) == 0 {
			return nil, &domainErrors.ErrBadRequest{Message: "skor pertandingan belum di update"}
		}
		for _, s := range match.Sets {
			if s.TeamAScore == 0 && s.TeamBScore == 0 {
				return nil, &domainErrors.ErrBadRequest{Message: "skor pertandingan belum di update"}
			}
		}
	}

	match.IsComplete = true
	match.IsCompletedBy = userUUID
	now := time.Now()
	match.IsCompletedAt = &now

	updatedUUID, err := s.matchRepo.Update(ctx, match)
	if err != nil {
		return nil, fmt.Errorf("s.matchRepo.Update(ctx, match) : %w", err)
	}

	return updatedUUID, nil
}

func (s *Service) ConfirmCompleteMatch(ctx context.Context, userUUID, matchUUID uuid.UUID) (*uuid.UUID, error) {
	match, err := s.matchRepo.GetByUUID(ctx, matchUUID)
	if err != nil {
		return nil, fmt.Errorf("s.matchRepo.GetByUUID(ctx, matchUUID=%s) : %w", matchUUID, err)
	}

	if match.TeamA.CaptainUserUUID != userUUID && match.TeamB.CaptainUserUUID != userUUID {
		return nil, &domainErrors.ErrUnauthorized{Message: "not authorized to complete match"}
	}

	if !match.IsComplete {
		return nil, &domainErrors.ErrBadRequest{Message: "match not yet complete"}
	}

	match.IsCompleteConfirmed = true
	match.IsCompleteConfirmedBy = userUUID
	now := time.Now()
	match.IsCompleteConfirmedAt = &now
	match.Status = string(enums.MATCH_STATUS_COMPLETED)

	updatedUUID, err := s.matchRepo.Update(ctx, match)
	if err != nil {
		return nil, fmt.Errorf("s.matchRepo.Update(ctx, match) : %w", err)
	}

	return updatedUUID, nil
}
