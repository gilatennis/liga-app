package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces"
	iCommon "liga-app/interfaces/common"

	"github.com/google/uuid"
)

type IGormRepo interface {
	Save(ctx context.Context, entity interfaces.IEntity, user *entities.User) (*uuid.UUID, error)
	Delete(ctx context.Context, entity interfaces.IEntity, user *entities.User) (*uuid.UUID, error)
	Get(
		ctx context.Context,
		p *iCommon.Pagination,
		model any,
		where *map[string]interface{},
	) (any, int64, error)
}
