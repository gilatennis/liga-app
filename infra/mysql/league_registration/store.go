package league_registration

import (
	"context"
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

func (s *Store) GetByUUID(ctx context.Context, leagueRegUUID uuid.UUID) (*entities.LeagueRegistration, error) {
	var leagueReg entities.LeagueRegistration
	err := s.db.WithContext(ctx).
		Where("uuid = ?", leagueRegUUID).
		First(&leagueReg).
		Error
	if err != nil {
		return nil, err
	}

	return &leagueReg, nil
}

func (s *Store) Update(ctx context.Context, leagueReg *entities.LeagueRegistration) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Save(leagueReg).
		Error
	if err != nil {
		return nil, err
	}

	return &leagueReg.UUID, nil
}
