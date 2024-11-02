package services

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"time"

	"github.com/google/uuid"
)

type ICourtService interface {
	Get(c context.Context, p *common.Pagination) ([]entities.Court, error)
	GetByUUID(c context.Context, courtUUID uuid.UUID) (*entities.Court, error)
	GetCourtAvailableDays(c context.Context, courtUUID uuid.UUID) ([]time.Time, error)
	GetAdminCourtAvailableDays(c context.Context, courtUUID uuid.UUID) ([]time.Time, error)
	GetCourtTimeSlotsByDay(c context.Context, courtUUID uuid.UUID, day time.Time) ([]entities.CourtTimeSlot, error)
}
