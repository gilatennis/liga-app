package auth

import (
	"context"
	"errors"
	"fmt"
	domainErrors "liga-app/domain/errors"
	"liga-app/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db,
	}
}

func (s *Store) CreateMagicLinkRegistration(ctx context.Context, reg *entities.MagicLinkRegistration) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Save(reg).
		Error
	if err != nil {
		return nil, err
	}

	return &reg.UUID, nil
}

func (s *Store) UpdateMagicLinkRegistration(ctx context.Context, reg *entities.MagicLinkRegistration) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Save(reg).
		Error
	if err != nil {
		return nil, err
	}

	return &reg.UUID, nil
}

func (s *Store) GetMagicLinkRegistrationByEmail(ctx context.Context, email string) (*entities.MagicLinkRegistration, error) {
	reg := entities.MagicLinkRegistration{}
	err := s.db.Where("email = ?", email).First(&reg).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &domainErrors.ErrNotFound{
			Message: fmt.Sprintf("s.db.Where(\"email = ?\", %s).First(&reg).Error: %v\n", email, err.Error()),
		}
	} else if err != nil {
		return nil, err
	}

	return &reg, nil
}

func (s *Store) CreateMagicLinkLogin(ctx context.Context, login *entities.MagicLinkLogin) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Save(login).
		Error
	if err != nil {
		return nil, err
	}

	return &login.UUID, nil
}

func (s *Store) UpdateMagicLinkLogin(ctx context.Context, login *entities.MagicLinkLogin) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Save(login).
		Error
	if err != nil {
		return nil, err
	}

	return &login.UUID, nil
}

func (s *Store) GetMagicLinkLoginByToken(ctx context.Context, token string) (*entities.MagicLinkLogin, error) {
	login := entities.MagicLinkLogin{}
	err := s.db.Where("token = ?", token).First(&login).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &domainErrors.ErrNotFound{
			Message: fmt.Sprintf("s.db.Where(\"token = ?\", %s).First(&reg).Error: %v\n", token, err.Error()),
		}
	} else if err != nil {
		return nil, err
	}

	return &login, nil
}
