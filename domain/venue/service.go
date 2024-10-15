package venue

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type Service struct {
	venueRepo repositories.IVenueRepository
}

func NewService(
	venueRepo repositories.IVenueRepository,
) *Service {
	return &Service{
		venueRepo: venueRepo,
	}
}

func (s *Service) Get(c context.Context, p *common.Pagination, f *repositories.VenueFilter) ([]entities.Venue, int64, error) {
	venues, count, err := s.venueRepo.Get(c, p, f)
	if err != nil {
		return nil, 0, fmt.Errorf("s.venueRepo.Get(c, p, f) : %w", err)
	}

	return venues, count, nil
}

func (s *Service) GetByUUID(c context.Context, venueUUID uuid.UUID) (*entities.Venue, error) {
	venue, err := s.venueRepo.GetByUUID(c, venueUUID)
	if err != nil {
		return nil, fmt.Errorf("s.venueRepo.GetByUUID(c, venueUUID=%s): %w", venueUUID, err)
	}

	return venue, nil
}
