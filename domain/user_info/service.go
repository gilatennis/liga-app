package user_info

import (
	"context"
	"liga-app/config"
	"liga-app/entities"
	"liga-app/interfaces"
	iCommon "liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"
	"mime/multipart"

	"github.com/google/uuid"
)

type Service struct {
	cfg         *config.Config
	wilayahRepo iRepo.IWilayahRepository
	userRepo    iRepo.IUserRepository
	fileStorage interfaces.IFileStorage
}

func NewService(
	cfg *config.Config,
	wilayahRepo iRepo.IWilayahRepository,
	userRepo iRepo.IUserRepository,
	fileStorage interfaces.IFileStorage,
) *Service {
	return &Service{
		cfg:         cfg,
		wilayahRepo: wilayahRepo,
		userRepo:    userRepo,
		fileStorage: fileStorage,
	}
}

func (s *Service) GetProvinceOptions(ctx context.Context) ([]entities.Wilayah, error) {
	p := iCommon.NewPagination(1, 1000, []iCommon.OrderBy{{Column: "name", Direction: "asc"}})
	f := iRepo.WilayahFilter{Type: 1}

	cities, err := s.wilayahRepo.Get(ctx, p, &f)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func (s *Service) SaveUserInfo(ctx context.Context, user *entities.User, profilePicture *multipart.File) (*uuid.UUID, error) {
	if profilePicture != nil {
		profilePicName, err := s.fileStorage.CreateFileName(ctx, *profilePicture)
		if err != nil {
			return nil, err
		}

		err = s.fileStorage.SaveFile(ctx, s.cfg.Minio.Bucket, profilePicName, *profilePicture)
		if err != nil {
			return nil, err
		}

		user.ProfilePicture = profilePicName
	}

	userUUID, err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return userUUID, nil
}
