package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
)

type VenueFilter struct{}

type IVenueRepository interface {
	Get(c context.Context, p *common.Pagination, f *VenueFilter) ([]entities.Venue, int64, error)
	GetByUUID(c context.Context, venueUUID uuid.UUID) (*entities.Venue, error)
}
