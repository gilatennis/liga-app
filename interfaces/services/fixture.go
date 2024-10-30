package services

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type IFixtureService interface {
	Get(ctx context.Context, p *common.Pagination, f *repositories.FixtureFilter) ([]entities.Fixture, int64, error)
	GetByUUID(ctx context.Context, fixtureUUID uuid.UUID) (*entities.Fixture, error)
}
