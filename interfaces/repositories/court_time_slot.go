package repositories

import (
	"context"
	"liga-app/entities"
)

type ICourtTimeSlotRepository interface {
	GetByCourtSlug(ctx context.Context, slug string) ([]entities.CourtTimeSlot, error)
}
