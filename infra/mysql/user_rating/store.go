package user_rating

import (
	"context"
	"fmt"
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

func (s *Store) Save(ctx context.Context, userRating *entities.UserRating) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Save(userRating).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db...Save(userRating) : %w", err)
	}

	return &userRating.UUID, nil
}
