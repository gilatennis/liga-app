package services

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type ILeagueService interface {
	Get(ctx context.Context, p *common.Pagination, f *repositories.LeagueFilter) ([]entities.League, int64, error)
	GetByUUID(ctx context.Context, leagueUUID uuid.UUID) (*entities.League, error)
	Join(ctx context.Context, leagueUUID uuid.UUID, user *entities.User) (string, error)
	CompleteLeagueRegistrationPayment(
		ctx context.Context,
		leagueRegUUID uuid.UUID,
		paymentID string,
		paymentAmount float64,
	) (*uuid.UUID, error)
	RegisterTeamToLeague(ctx context.Context, leagueReg *entities.LeagueRegistration, teamName string, playerUUIDs []uuid.UUID) (*uuid.UUID, error)
}
