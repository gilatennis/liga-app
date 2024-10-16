package services

import (
	"context"
	"liga-app/entities"
	"time"

	"github.com/google/uuid"
)

type ICourtService interface {
	GetByUUID(c context.Context, courtUUID uuid.UUID) (*entities.Court, error)
	GetCourtAvailableDays(c context.Context, courtUUID uuid.UUID) ([]time.Time, error)
	GetCourtTimeSlotsByDay(c context.Context, courtUUID uuid.UUID, day time.Time) ([]entities.CourtTimeSlot, error)
}
