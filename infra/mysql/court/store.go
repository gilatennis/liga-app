package court

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db}
}

func (s *Store) Get(c context.Context, p *common.Pagination) ([]entities.Court, error) {
	courts := []entities.Court{}
	err := s.db.Find(&courts).Error
	if err != nil {
		return nil, fmt.Errorf("s.db.Find(&courts): %w", err)
	}

	return courts, nil
}

func (s *Store) GetByUUID(c context.Context, courtUUID uuid.UUID) (*entities.Court, error) {
	court := entities.Court{}

	err := s.db.WithContext(c).First(&court, "uuid = ?", courtUUID).Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(c).First(&court, \"uuid = ?\", courtUUID).Error: %w", err)
	}

	return &court, nil
}
