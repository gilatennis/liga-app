package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
)

type ICourtRepository interface {
	Get(c context.Context, p *common.Pagination) ([]entities.Court, error)
	GetByUUID(c context.Context, courtUUID uuid.UUID) (*entities.Court, error)
}
