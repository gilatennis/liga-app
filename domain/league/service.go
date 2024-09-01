package league

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
)

type Service struct {
	leagueRepo repositories.ILeagueRepository
}

func NewService(leageRepo repositories.ILeagueRepository) *Service {
	return &Service{
		leagueRepo: leageRepo,
	}
}

func (s *Service) Get(ctx context.Context, p *common.Pagination, f *repositories.LeagueFilter) ([]entities.League, int64, error) {
	leagues, count, err := s.leagueRepo.Get(ctx, p, f)
	if err != nil {
		return nil, 0, err
	}

	return leagues, count, err
}

func (s *Service) Join(ctx context.Context, user *entities.User) error {
	// TODO:
	return nil
}
