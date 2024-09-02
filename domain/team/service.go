package team

import (
	"context"
	"liga-app/entities"
	iCommon "liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

type Service struct {
	teamRepo iRepo.ITeamRepository
}

func NewService(teamRepo iRepo.ITeamRepository) *Service {
	return &Service{teamRepo: teamRepo}
}

func (s *Service) GetByUUID(ctx context.Context, teamUUID uuid.UUID) (*entities.Team, error) {
	team, err := s.teamRepo.GetByUUID(ctx, teamUUID)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (s *Service) GetUserTeamInLeague(ctx context.Context, userUUID uuid.UUID, leagueUUID uuid.UUID) (*entities.Team, error) {
	p := iCommon.NewPagination(1, 1, []iCommon.OrderBy{{Column: "created_at", Direction: "desc"}})
	f := &iRepo.TeamFilter{LeagueUUID: &leagueUUID, PlayerUUID: &userUUID}
	teams, _, err := s.teamRepo.Get(ctx, p, f)
	if err != nil {
		return nil, err
	}
	pp.Println("teams", teams)

	if len(teams) == 0 {
		return nil, nil
	}

	return &teams[0], nil
}
