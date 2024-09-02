package team

import (
	"context"
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
	return &Store{db}
}

func (s *Store) Get(ctx context.Context, p *common.Pagination, f *repositories.TeamFilter) ([]entities.Team, int64, error) {
	teams := []entities.Team{}
	var count int64
	db := s.db.WithContext(ctx).
		Preload("Captain").
		Preload("Players").
		Offset(p.Offset).
		Limit(p.Limit).
		Joins("JOIN league_teams ON league_teams.team_uuid = teams.uuid").
		Joins("JOIN team_players ON team_players.team_uuid = teams.uuid")

	if f != nil {
		if f.LeagueUUID != nil {
			db.Where("league_teams.league_uuid = ?", *f.LeagueUUID)
		}
		if f.PlayerUUID != nil {
			db.Where("team_players.player_uuid = ?", *f.PlayerUUID)
		}
	}

	err := db.Find(&teams).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return teams, count, nil
}

func (s *Store) GetByUUID(ctx context.Context, teamUUID uuid.UUID) (*entities.Team, error) {
	team := entities.Team{}
	err := s.db.WithContext(ctx).
		Where("uuid = ?", teamUUID).
		Preload("Captain").
		Preload("Players").
		First(&team).Error
	if err != nil {
		return nil, err
	}

	return &team, nil
}
