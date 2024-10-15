package services

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type IVenueService interface {
	Get(c context.Context, p *common.Pagination, f *repositories.VenueFilter) ([]entities.Venue, int64, error)
	GetByUUID(c context.Context, venueUUID uuid.UUID) (*entities.Venue, error)
}
