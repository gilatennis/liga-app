package admin

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
	return &Store{db}
}

func (s *Store) GetByUserUUID(ctx context.Context, userUUID uuid.UUID) (*entities.Admin, error) {
	adm := entities.Admin{}

	err := s.db.Where("user_uuid = ?", userUUID).First(&adm).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &domainErrors.ErrNotFound{Message: "not found"}
	} else if err != nil {
		return nil, fmt.Errorf("s.db.Where(user_uuid = ?, userUUID=%s): %w", userUUID, err)
	}

	return &adm, nil
}
