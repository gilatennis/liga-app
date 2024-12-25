package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
)

type FixtureFilter struct {
	LeagueUUID *uuid.UUID
	Format     *enums.FIXTURE_FORMAT
}

type IFixtureRepository interface {
	Get(ctx context.Context, p *common.Pagination, f *FixtureFilter) ([]entities.Fixture, int64, error)
	GetByUUID(ctx context.Context, fixtureUUID uuid.UUID) (*entities.Fixture, error)
}
