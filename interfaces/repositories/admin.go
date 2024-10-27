package repositories

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type IAdminRepository interface {
	GetByUserUUID(ctx context.Context, userUUID uuid.UUID) (*entities.Admin, error)
}
