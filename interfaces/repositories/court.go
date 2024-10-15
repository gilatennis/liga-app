package repositories

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type ICourtRepository interface {
	GetByUUID(c context.Context, courtUUID uuid.UUID) (*entities.Court, error)
}
