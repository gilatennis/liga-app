package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"time"

	"github.com/google/uuid"
)

type CourtTimeSlotFilter struct {
	CourtTimeSlotUUIDs []uuid.UUID
}

type ICourtTimeSlotRepository interface {
	GetByCourtSlug(ctx context.Context, slug string, rangeStart time.Time, rangeEnd time.Time) ([]entities.CourtTimeSlot, error)
	Get(ctx context.Context, p *common.Pagination, f *CourtTimeSlotFilter) ([]entities.CourtTimeSlot, int64, error)
}
