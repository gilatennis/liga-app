package user_info

import (
	"context"
	"liga-app/entities"
	iCommon "liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type Service struct {
	wilayahRepo iRepo.IWilayahRepository
	userRepo    iRepo.IUserRepository
}

func NewService(
	wilayahRepo iRepo.IWilayahRepository,
	userRepo iRepo.IUserRepository,
) *Service {
	return &Service{
		wilayahRepo: wilayahRepo,
		userRepo:    userRepo,
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

func (s *Service) SaveUserInfo(ctx context.Context, user *entities.User) (*uuid.UUID, error) {
	userUUID, err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return userUUID, nil
}
