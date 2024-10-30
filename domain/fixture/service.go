package fixture

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type Service struct {
	fixtureRepo repositories.IFixtureRepository
}

func NewService(fixtureRepo repositories.IFixtureRepository) *Service {
	return &Service{
		fixtureRepo: fixtureRepo,
	}
}

func (s *Service) Get(ctx context.Context, p *common.Pagination, f *repositories.FixtureFilter) ([]entities.Fixture, int64, error) {
	fixtures, count, err := s.fixtureRepo.Get(ctx, p, f)
	if err != nil {
		return nil, 0, fmt.Errorf("s.fixtureRepo.Get(ctx, p, f): %w", err)
	}

	return fixtures, count, nil
}

func (s *Service) GetByUUID(ctx context.Context, fixtureUUID uuid.UUID) (*entities.Fixture, error) {
	fixture, err := s.fixtureRepo.GetByUUID(ctx, fixtureUUID)
	if err != nil {
		return nil, fmt.Errorf("s.fixtureRepo.GetByUUID(ctx, fixtureUUID=%s): %w", fixtureUUID, err)
	}

	return fixture, nil
}
