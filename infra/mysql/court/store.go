package court

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

func (s *Store) GetByUUID(c context.Context, courtUUID uuid.UUID) (*entities.Court, error) {
	court := entities.Court{}

	err := s.db.WithContext(c).First(&court, "uuid = ?", courtUUID).Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(c).First(&court, \"uuid = ?\", courtUUID).Error: %w", err)
	}

	return &court, nil
}
