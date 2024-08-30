package repositories

import (
	"liga-app/entities"
	"liga-app/interfaces"

	"github.com/google/uuid"
)

type IGormRepo interface {
	Save(entity interfaces.IEntity, user *entities.User) (*uuid.UUID, error)
	Delete(entity interfaces.IEntity, user *entities.User) (*uuid.UUID, error)
	// Get(entity interfaces.IEntity)
}
