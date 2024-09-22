package services

import (
	"context"
	"liga-app/entities"
	"mime/multipart"

	"github.com/google/uuid"
)

type IUserInfoService interface {
	GetProvinceOptions(ctx context.Context) ([]entities.Wilayah, error)
	SaveUserInfo(ctx context.Context, user *entities.User, profilePicture *multipart.File) (*uuid.UUID, error)
}
