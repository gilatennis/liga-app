package invite

import (
	"context"
	domainErrors "liga-app/domain/errors"
	"liga-app/enums"
	iRepo "liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type Service struct {
	leagueRegRepo iRepo.ILeagueRegistrationRepository
	leagueRepo    iRepo.ILeagueRepository
	teamRepo      iRepo.ITeamRepository
}

func NewService(
	leagueRegRepo iRepo.ILeagueRegistrationRepository,
	leagueRepo iRepo.ILeagueRepository,
	teamRepo iRepo.ITeamRepository,
) *Service {
	return &Service{
		leagueRegRepo: leagueRegRepo,
		leagueRepo:    leagueRepo,
		teamRepo:      teamRepo,
	}
}

func (s *Service) AcceptLeagueTeamInvite(ctx context.Context, userUUID, leagueUUID, teamUUID uuid.UUID) error {
	// if user already in team, can not join again.
	team, err := s.teamRepo.GetByUUID(ctx, teamUUID)
	if err != nil {
		return err
	}

	isAlreadyInTeam := false
	for _, p := range team.Players {
		if p.UUID == userUUID {
			isAlreadyInTeam = true
		}
	}

	if isAlreadyInTeam {
		return &domainErrors.ErrConflict{Message: "user sudah masuk ke dalam tim"}
	}

	league, err := s.leagueRepo.GetByUUID(ctx, leagueUUID)
	if err != nil {
		return err
	}

	if league.ParticipationCategory == enums.CATEGORY_PARTICIPATION_TEAMS {
		if len(team.Players) == 9 {
			return &domainErrors.ErrConflict{Message: "jumlah anggota tim sudah mencapai batas maksimal"}
		}
	}

	if league.ParticipationCategory == enums.CATEGORY_PARTICIPATION_DOUBLES {
		if len(team.Players) == 2 {
			return &domainErrors.ErrConflict{Message: "jumlah anggota tim sudah mencapai batas maksimal"}
		}
	}

	err = s.leagueRegRepo.JoinExistingTeamInLeague(ctx, userUUID, leagueUUID, teamUUID)
	if err != nil {
		return err
	}
	return nil
}
