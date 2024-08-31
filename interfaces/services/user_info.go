package services

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type IUserInfoService interface {
	GetProvinceOptions(ctx context.Context) ([]entities.Wilayah, error)
	SaveUserInfo(ctx context.Context, user *entities.User) (*uuid.UUID, error)
}
