package services

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

type IUserService interface {
	CreateUser(ctx context.Context, userUUID uuid.UUID, e *entities.User) (uuid.UUID, error)
	GetUsers(ctx context.Context, p *common.Pagination, f *GetUserListFilter) (*[]entities.User, int64, error)
	GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (*entities.User, error)
	UpdateUser(ctx context.Context, userUUID uuid.UUID, e *entities.User) (uuid.UUID, error)
	DeleteUser(ctx context.Context, actingUserUUID uuid.UUID, userUUID uuid.UUID) error
}
