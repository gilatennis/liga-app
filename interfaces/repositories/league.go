package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
)

type LeagueFilter struct {
	CityID     *string
	ProvinceID *string
	UserUUID   *uuid.UUID
}

type ILeagueRepository interface {
	Get(ctx context.Context, p *common.Pagination, f *LeagueFilter) ([]entities.League, int64, error)
	GetByUUID(ctx context.Context, leagueUUID uuid.UUID) (*entities.League, error)
	Create(ctx context.Context, e *entities.League) (*uuid.UUID, error)
	Update(ctx context.Context, e *entities.League) (*uuid.UUID, error)
	Delete(ctx context.Context, e *entities.League) (*uuid.UUID, error)
}
