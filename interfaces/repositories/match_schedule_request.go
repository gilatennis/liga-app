package repositories

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type IMatchScheduleRequestRepository interface {
	Create(ctx context.Context, req *entities.MatchScheduleRequest) (*uuid.UUID, error)
	GetByUUID(ctx context.Context, reqUUID uuid.UUID) (*entities.MatchScheduleRequest, error)
	Update(ctx context.Context, req *entities.MatchScheduleRequest) (*uuid.UUID, error)
	AcceptMatchScheduleRequest(ctx context.Context, req *entities.MatchScheduleRequest) (*uuid.UUID, error)
}
