package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
)

type GetUserListFilter struct {
	Name        *string
	Email       *string
	LoginMethod *string
}

type IUserRepository interface {
	CreateUser(ctx context.Context, t *entities.User) (*uuid.UUID, error)
	GetUsers(
		ctx context.Context,
		p *common.Pagination,
		f *GetUserListFilter,
	) ([]entities.User, int64, error)
	GetUserByUUID(ctx context.Context, tUUID uuid.UUID) (*entities.User, error)
	UpdateUser(ctx context.Context, t *entities.User) (*uuid.UUID, error)
	DeleteUser(ctx context.Context, userUUID uuid.UUID) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}
