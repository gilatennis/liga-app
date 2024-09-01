package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces"

	"github.com/google/uuid"
)

type IGormRepo interface {
	Save(ctx context.Context, entity interfaces.IEntity, user *entities.User) (*uuid.UUID, error)
	Delete(ctx context.Context, entity interfaces.IEntity, user *entities.User) (*uuid.UUID, error)
	// Get(entity interfaces.IEntity)
}
