package match_set

import (
	"context"
	"errors"
	"fmt"
	"liga-app/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db}
}

func (s *Store) SaveMatchSet(ctx context.Context, userUUID uuid.UUID, matchSet *entities.MatchSet) (*uuid.UUID, error) {
	existingMatchSet := entities.MatchSet{}

	tx := s.db.WithContext(ctx).Begin()

	err := tx.
		Where("match_uuid = ? AND number = ?", matchSet.MatchUUID, matchSet.Number).
		First(&existingMatchSet).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("s.db.Where(...).First(&existingMatchSet).Error : %w", err)
	}

	// Create
	if errors.Is(err, gorm.ErrRecordNotFound) {
		matchSet.CreatedBy = userUUID
		err = tx.Create(&matchSet).Error
		if err != nil {
			return nil, fmt.Errorf("tx.Create(&matchSet).Error : %w", err)
		}

		err = tx.Commit().Error
		if err != nil {
			return nil, fmt.Errorf("tx.Commit().Error")
		}

		return &matchSet.UUID, nil
	}

	// Update
	existingMatchSet.UpdatedBy = userUUID
	existingMatchSet.TeamAScore = matchSet.TeamAScore
	existingMatchSet.TeamBScore = matchSet.TeamBScore
	err = tx.Save(&existingMatchSet).Error
	if err != nil {
		return nil, fmt.Errorf("tx.Save(&existingMatchSet).Error : %w", err)
	}

	match := entities.Match{}
	err = tx.First(&match, "uuid = ?", matchSet.MatchUUID).Error
	if err != nil {
		return nil, fmt.Errorf("tx.First(&match, \"uuid = ?\", matchSet.MatchUUID).Error")
	}

	match.UpdatedAt = time.Now()
	err = tx.Save(&match).Error
	if err != nil {
		return nil, fmt.Errorf("tx.Save(&match).Error")
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, fmt.Errorf("tx.Commit().Error")
	}

	return &existingMatchSet.UUID, nil
}
