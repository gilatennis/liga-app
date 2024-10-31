package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
)

type TeamFilter struct {
	LeagueUUID *uuid.UUID
	PlayerUUID *uuid.UUID
	GroupUUID  *uuid.UUID
	NotInGroup *bool
}

type ITeamRepository interface {
	Get(ctx context.Context, p *common.Pagination, f *TeamFilter) ([]entities.Team, int64, error)
	GetByUUID(ctx context.Context, teamUUID uuid.UUID) (*entities.Team, error)
	Update(ctx context.Context, team *entities.Team) (*uuid.UUID, error)
}
