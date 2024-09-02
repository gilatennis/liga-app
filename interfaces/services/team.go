package services

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type ITeamService interface {
	GetUserTeamInLeague(ctx context.Context, userUUID uuid.UUID, leagueUUID uuid.UUID) (*entities.Team, error)
	GetByUUID(ctx context.Context, teamUUID uuid.UUID) (*entities.Team, error)
}
