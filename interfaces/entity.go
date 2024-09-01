package interfaces

import (
	"liga-app/entities"

	"github.com/google/uuid"
)

type IEntity interface {
	GetUUID() *uuid.UUID
	SetCreatedBy(user *entities.User)
	SetUpdatedBy(user *entities.User)
	SetDeletedBy(user *entities.User)
}
