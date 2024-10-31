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
	groupRepo                iRepo.IGroupRepository
}

func NewService(
	matchRepo iRepo.IMatchRepository,
	courtTimeSlotRepo iRepo.ICourtTimeSlotRepository,
	matchScheduleRequestRepo iRepo.IMatchScheduleRequestRepository,
	matchSetRepo iRepo.IMatchSetRepository,
	groupRepo iRepo.IGroupRepository,
) *Service {
	return &Service{
		matchRepo:                matchRepo,
		courtTimeSlotRepo:        courtTimeSlotRepo,
		matchScheduleRequestRepo: matchScheduleRequestRepo,
		matchSetRepo:             matchSetRepo,
		groupRepo:                groupRepo,
	}
}

func (s *Service) Get(ctx context.Context, p *common.Pagination, f *iService.MatchFilter) ([]entities.Match, int64, error) {
	repoFilter := &iRepo.MatchFilter{
		UserUUID:   f.UserUUID,
		GroupUUID:  f.GroupUUID,
		LeagueUUID: f.LeagueUUID,
	}

	// INFO: Rounds start on wednesday
	startOfWeek := s.getStartOfWeek()
	wednesday := startOfWeek.Add(2 * 24 * time.Hour)
	nextWeekTuesday := wednesday.Add(6 * 24 * time.Hour).Add(23 * time.Hour).Add(59 * time.Minute).Add(59 * time.Second)

	if f != nil {
		if f.Status != nil {
			switch *f.Status {
			case iService.MATCH_FILTER_STATUS_CURRENT_ROUND:
				repoFilter.RoundIsAfter = &wednesday
				repoFilter.RoundIsBefore = &nextWeekTuesday

			case iService.MATCH_FILTER_STATUS_UPCOMING:
				repoFilter.RoundIsAfter = &wednesday

			case iService.MATCH_FILTER_STATUS_COMPLETED:
				status := enums.MATCH_STATUS_COMPLETED
				repoFilter.Status = &status

			case iService.MATCH_FILTER_STATUS_OVERDUE:
				status := enums.MATCH_STATUS_WAITING
				repoFilter.Status = &status
				repoFilter.RoundIsBefore = &wednesday

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

	// INFO: rounds start on wednesday
	startOfWeek := s.getStartOfWeek().Add(2 * 24 * time.Hour)
	endOfWeek := startOfWeek.Add(6 * 24 * time.Hour).Add(23 * time.Hour).Add(59 * time.Minute).Add(59 * time.Second)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	/*
		fmt.Println("startOfWeek", startOfWeek)
		fmt.Println("endOfWeek", endOfWeek)
		fmt.Println("m.Round.Start", m.Round.Start.In(loc))
		fmt.Println("m.Round.End", m.Round.End.In(loc))
		fmt.Println("m.Round.End.Before(endOfWeek)", m.Round.End.Before(endOfWeek))
	*/

	if m.Round.Start != nil {
		localTime := m.Round.Start.In(loc)
		m.Round.Start = &localTime
	}

	if m.Round.End != nil {
		localTime := m.Round.End.In(loc)
		m.Round.End = &localTime
	}

	if m.Round.Start != nil && m.Round.Start.Equal(startOfWeek) && m.Round.End != nil && m.Round.End.Equal(endOfWeek) {
		m.FilterStatus = string(iService.MATCH_FILTER_STATUS_CURRENT_ROUND)
	} else if m.Round.End != nil && startOfWeek.Before(*m.Round.End) && m.Status == string(enums.MATCH_STATUS_WAITING) {
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

func (s *Service) GenerateRoundRobinMatches(ctx context.Context, groupUUID uuid.UUID, isHomeAndAway bool) error {
	g, err := s.groupRepo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return fmt.Errorf("s.groupRepo.GetByUUID(ctx, groupUUID=%s): %w", groupUUID, err)
	}

	teamCount := len(g.Teams)

	teamIDs := make([]string, teamCount)
	for i, team := range g.Teams {
		teamIDs[i] = team.UUID.String()
	}

	if len(teamIDs) < teamCount {
		for i := len(teamIDs); i < teamCount; i++ {
			teamIDs = append(teamIDs, "TBA")
		}
	}

	// assign number to teamIDs
	numToTeamID := make(map[int]string)
	var teamNums []int
	const BYE = "BYE"

	if len(teamIDs)%2 != 0 {
		teamNums = make([]int, len(teamIDs)+1)
		teamIDs = append(teamIDs, BYE)
	} else {
		teamNums = make([]int, len(teamIDs))
	}

	for i, teamID := range teamIDs {
		numToTeamID[i] = teamID
		teamNums[i] = i
	}

	numRounds := len(teamIDs) - 1

	if isHomeAndAway {
		numRounds = numRounds * 2
	}
	fmt.Println("numRounds", numRounds)

	rounds := []entities.Round{}
	for i := 0; i < numRounds; i++ {
		r := entities.Round{
			Base: entities.Base{
				UUID: uuid.New(),
			},
			GroupUUID: groupUUID,
			Number:    i + 1,
		}

		for j := 0; j < len(teamNums)/2; j++ {
			left := teamNums[j]
			right := teamNums[len(teamNums)-j-1]

			opp1ID := numToTeamID[left]
			opp2ID := numToTeamID[right]

			match := entities.Match{
				Base: entities.Base{
					UUID: uuid.New(),
				},
				Number:    j + 1,
				RoundUUID: r.UUID,
				// RoundNumber: r.Number,
				Status: string(enums.MATCH_STATUS_WAITING),
			}

			oppIDs := []string{opp1ID, opp2ID}
			for i := 0; i < 2; i++ {
				oppID := oppIDs[i]

				if i == 0 {
					if oppID == BYE {
						match.IsByeTeamA = true
						continue
					}

					if oppID == "TBA" {
						continue
					}

					oppUUID, err := uuid.Parse(oppID)
					if err != nil {
						err = fmt.Errorf("uuid.Parse(oppID=%s): %w", oppID, err)
						return err
					}

					match.TeamAUUID = &oppUUID
				}

				if i == 1 {
					if oppID == BYE {
						match.IsByeTeamB = true
						continue
					}

					if oppID == "TBA" {
						continue
					}

					oppUUID, err := uuid.Parse(oppID)
					if err != nil {
						err = fmt.Errorf("uuid.Parse(oppID=%s): %w", oppID, err)
						return err
					}

					match.TeamBUUID = &oppUUID
				}
			}

			r.Matches = append(r.Matches, match)
		}

		rounds = append(rounds, r)

		if i == numRounds-1 {
			break
		}
		// rotate teamIDs
		rTemp := 0
		for k := 0; k < len(teamNums); k++ {
			// keep first idx
			if k == 0 {
				continue
			}
			if k == 1 {
				rTemp = teamNums[k]
				teamNums[k] = teamNums[len(teamNums)-1]
				continue
			}
			temp := teamNums[k]
			teamNums[k] = rTemp
			rTemp = temp
		}
	}

	err = s.matchRepo.CreateGeneratedMatches(ctx, rounds, g.Fixture.LeagueUUID, g.UUID)
	if err != nil {
		return fmt.Errorf("s.matchRepo.CreateGeneratedMatches(ctx, rounds, g.Fixture.LeagueUUID, g.UUID) : %w", err)
	}

	return nil
}
