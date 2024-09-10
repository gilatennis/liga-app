package services

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type ITeamService interface {
	Get(ctx context.Context, p *common.Pagination, f *repositories.TeamFilter) ([]entities.Team, int64, error)
	GetUserTeamInLeague(ctx context.Context, userUUID uuid.UUID, leagueUUID uuid.UUID) (*entities.Team, error)
	GetByUUID(ctx context.Context, teamUUID uuid.UUID) (*entities.Team, error)
	GetLeagueTeamInviteLink(ctx context.Context, leagueUUID uuid.UUID, teamUUID uuid.UUID) string
}
