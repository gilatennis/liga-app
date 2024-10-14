package repositories

import (
	"context"
	"liga-app/entities"
	"time"
)

type ICourtTimeSlotRepository interface {
	GetByCourtSlug(ctx context.Context, slug string, rangeStart time.Time, rangeEnd time.Time) ([]entities.CourtTimeSlot, error)
}
