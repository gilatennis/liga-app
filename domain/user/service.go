package user

import (
	"context"
	"errors"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type Service struct {
	userRepo repositories.IUserRepository
}

func NewService(userRepo repositories.IUserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) CreateUser(ctx context.Context, userUUID uuid.UUID, e *entities.User) (*uuid.UUID, error) {
	return nil, errors.New("unimplemented")
}

func (s *Service) GetUsers(ctx context.Context, p *common.Pagination, f *repositories.GetUserListFilter) ([]entities.User, int64, error) {
	users, count, err := s.userRepo.GetUsers(ctx, p, f)
	if err != nil {
		return nil, 0, fmt.Errorf("s.userRepo.GetUsers(ctx, p, f): %w", err)
	}

	return users, count, err
}

func (s *Service) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.GetUserByUUID(ctx, userUUID=%s): %w", userUUID, err)
	}
	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.GetUserByEmail(ctx, email=%s): %w", email, err)
	}
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, userUUID uuid.UUID, e *entities.User) (*uuid.UUID, error) {
	updatedUserUUID, err := s.userRepo.UpdateUser(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.UpdateUser(ctx, e) : %w", err)
	}

	return updatedUserUUID, nil
}

func (s *Service) DeleteUser(ctx context.Context, actingUserUUID uuid.UUID, userUUID uuid.UUID) error {
	return errors.New("unimplemented")
}
