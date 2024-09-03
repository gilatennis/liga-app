package repositories

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type ILeagueRegistrationRepository interface {
	GetByUUID(ctx context.Context, leagueRegUUID uuid.UUID) (*entities.LeagueRegistration, error)
	Update(ctx context.Context, leagueReg *entities.LeagueRegistration) (*uuid.UUID, error)
	CompleteRegistration(ctx context.Context, leagueReg *entities.LeagueRegistration) (*uuid.UUID, error)
	JoinExistingTeamInLeague(ctx context.Context, userUUID, leagueUUID, teamUUID uuid.UUID) error
}
