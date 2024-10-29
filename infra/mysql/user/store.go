package user

import (
	"context"
	"errors"
	"fmt"
	domainErrors "liga-app/domain/errors"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

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

func (s *Store) CreateUser(ctx context.Context, t *entities.User) (*uuid.UUID, error) {
	err := s.db.Save(t).Error
	if err != nil {
		return nil, err
	}

	return &t.UUID, nil
}

func (s *Store) GetUsers(
	ctx context.Context,
	p *common.Pagination,
	f *repositories.GetUserListFilter,
) ([]entities.User, int64, error) {
	dbHandle := s.db.WithContext(ctx).
		Offset(p.Offset).
		Limit(p.Limit)

	users := []entities.User{}
	var count int64

	if f != nil {
		if f.LoginMethod != nil {
			dbHandle.Where("login_method = ?", *f.LoginMethod)
		}

		if f.Name != nil {
			dbHandle.Where("name = ?", *f.Name)
		}

		if f.Email != nil {
			dbHandle.Where("email = ?", *f.Email)
		}
	}

	for _, o := range p.OrderBy {
		dbHandle = dbHandle.Order(fmt.Sprintf("%s %s", o.Column, o.Direction))
	}

	err := dbHandle.Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	err = dbHandle.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (s *Store) GetUserByUUID(ctx context.Context, tUUID uuid.UUID) (*entities.User, error) {
	user := entities.User{}
	err := s.db.Where("uuid = ?", tUUID).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &domainErrors.ErrNotFound{
			Message: fmt.Sprintf("s.db.Where(\"uuid = ?\", %s).First(&user).Error: %v\n", tUUID, err.Error()),
		}
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) UpdateUser(ctx context.Context, t *entities.User) (*uuid.UUID, error) {
	err := s.db.Save(t).Error
	if err != nil {
		return nil, err
	}

	return &t.UUID, nil
}

func (s *Store) DeleteUser(ctx context.Context, userUUID uuid.UUID) error {
	err := s.db.Where("uuid = ?", userUUID).
		Delete(&entities.User{}).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := entities.User{}
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("a")
		return nil, &domainErrors.ErrNotFound{
			Message: fmt.Sprintf("s.db.Where(email = %s).First(&user): %v\n", email, err.Error()),
		}
	} else if err != nil {
		fmt.Println("b")
		return nil, err
	}

	return &user, nil
}
