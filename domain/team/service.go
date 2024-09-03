package team

import (
	"context"
	"fmt"
	"liga-app/config"
	"liga-app/entities"
	iCommon "liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type Service struct {
	cfg      *config.Config
	teamRepo iRepo.ITeamRepository
}

func NewService(
	cfg *config.Config,
	teamRepo iRepo.ITeamRepository,
) *Service {
	return &Service{
		cfg:      cfg,
		teamRepo: teamRepo,
	}
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

	if len(teams) == 0 {
		return nil, nil
	}

	return &teams[0], nil
}

func (s *Service) GetLeagueTeamInviteLink(ctx context.Context, leagueUUID uuid.UUID, teamUUID uuid.UUID) string {
	inviteLink := fmt.Sprintf("%s/invite?league_uuid=%s&team_uuid=%s", s.cfg.App.BaseURL, leagueUUID, teamUUID)
	return inviteLink
}
