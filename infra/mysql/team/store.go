package team

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	db := s.db.WithContext(ctx).Session(&gorm.Session{}).Debug().
		Distinct().
		Preload("Captain").
		Preload("Players").
		Preload("Players.Ratings").
		Offset(p.Offset).
		Limit(p.Limit).
		Joins("JOIN league_teams ON league_teams.team_uuid = teams.uuid").
		Joins("JOIN team_players ON team_players.team_uuid = teams.uuid").
		Joins("JOIN league_registrations lr ON lr.league_uuid = league_teams.league_uuid").
		Joins("LEFT JOIN group_teams gt ON gt.team_uuid = teams.uuid").
		Where("lr.payment_status = 'PAID'").
		Where("league_teams.deleted_at IS NULL")

	if f != nil {
		if f.LeagueUUID != nil {
			db = db.Where("league_teams.league_uuid = ?", *f.LeagueUUID)
		}
		if f.PlayerUUID != nil {
			db = db.Where("team_players.player_uuid = ?", *f.PlayerUUID)
		}
		if f.NotInGroup != nil && *f.NotInGroup {
			db = db.Where("(gt.group_uuid IS NULL)")
		}
	}

	err := db.Find(&teams).Error
	if err != nil {
		return nil, 0, err
	}

	dbCount := s.db.WithContext(ctx).Debug().
		Model(&entities.Team{}).
		Distinct("teams.uuid").
		Joins("JOIN league_teams ON league_teams.team_uuid = teams.uuid").
		Joins("JOIN team_players ON team_players.team_uuid = teams.uuid").
		Joins("JOIN league_registrations lr ON lr.league_uuid = league_teams.league_uuid").
		Joins("LEFT JOIN group_teams gt ON gt.team_uuid = teams.uuid").
		Where("lr.payment_status = 'PAID'").
		Where("league_teams.deleted_at IS NULL")

	if f != nil {
		if f.LeagueUUID != nil {
			dbCount = dbCount.Where("league_teams.league_uuid = ?", *f.LeagueUUID)
		}
		if f.PlayerUUID != nil {
			dbCount = dbCount.Where("team_players.player_uuid = ?", *f.PlayerUUID)
		}
		if f.NotInGroup != nil && *f.NotInGroup {
			dbCount = dbCount.Where("(gt.group_uuid IS NULL)")
		}
	}

	err = dbCount.Count(&count).Error
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
		Preload("Players.Ratings").
		First(&team).Error
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (s *Store) Update(ctx context.Context, team *entities.Team) (*uuid.UUID, error) {
	err := s.db.
		Omit(clause.Associations).
		Save(team).Error
	if err != nil {
		return nil, err
	}

	return &team.UUID, nil
}
