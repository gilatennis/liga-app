package services

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type IGroupService interface {
	GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*entities.Group, error)
	Create(ctx context.Context, group *entities.Group) (*uuid.UUID, error)
	AddTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error)
	RemoveTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error)
	GenerateGroupResults(ctx context.Context, groupUUID uuid.UUID) error
}
