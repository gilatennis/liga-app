package league

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

func (s *Store) Get(ctx context.Context, p *common.Pagination, f *repositories.LeagueFilter) ([]entities.League, int64, error) {
	leagues := []entities.League{}
	var count int64

	dbHandle := s.db.WithContext(ctx).
		Preload("Wilayah").
		Offset(p.Offset).
		Limit(p.Limit)

	if f != nil {
		if f.UserUUID != nil {
			dbHandle = dbHandle.
				Joins("JOIN league_players ON league_players.league_uuid = leagues.uuid").
				Where("league_players.player_uuid = ?", f.UserUUID)
		}
	}

	err := dbHandle.
		Find(&leagues).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}

	leagueUUIDs := []uuid.UUID{}
	for _, l := range leagues {
		_l := l
		leagueUUIDs = append(leagueUUIDs, _l.UUID)
	}

	type Temp struct {
		LeagueUUID  uuid.UUID
		TeamCount   int64
		PlayerCount int64
	}
	temps := []Temp{}

	err = s.db.WithContext(ctx).
		Raw(`
      select
        l.uuid AS league_uuid,
        COUNT(DISTINCT lt.uuid) AS team_count,
        COUNT(DISTINCT lp.uuid) AS player_count
      from leagues l
        left join league_teams lt on l.uuid = lt.league_uuid 
        left join league_players lp on l.uuid = lp.league_uuid
      WHERE l.uuid in ? AND lt.deleted_at IS NULL AND lp.deleted_at IS NULL
      GROUP BY l.uuid
    `, leagueUUIDs).
		Scan(&temps).
		Error
	if err != nil {
		return nil, 0, err
	}

	for i, l := range leagues {
		l.TeamCount = temps[i].TeamCount
		l.PlayerCount = temps[i].PlayerCount
		leagues[i] = l
	}

	return leagues, count, nil
}

func (s *Store) GetByUUID(ctx context.Context, leagueUUID uuid.UUID) (*entities.League, error) {
	league := entities.League{}
	err := s.db.Where("uuid = ?", leagueUUID).
		Preload("Wilayah").
		First(&league).
		Error
	if err != nil {
		return nil, err
	}

	type Temp struct {
		LeagueUUID  uuid.UUID
		TeamCount   int64
		PlayerCount int64
	}
	temp := Temp{}
	err = s.db.WithContext(ctx).
		Raw(`
      select
        l.uuid AS league_uuid,
        COUNT(DISTINCT lt.uuid) AS team_count,
        COUNT(DISTINCT lp.uuid) AS player_count
      from leagues l
        left join league_teams lt on l.uuid = lt.league_uuid 
        left join league_players lp on l.uuid = lp.league_uuid
      WHERE l.uuid = ? AND lt.deleted_at IS NULL AND lp.deleted_at IS NULL
      GROUP BY l.uuid
    `, league.UUID).
		Scan(&temp).
		Error
	if err != nil {
		return nil, err
	}

	league.TeamCount = temp.TeamCount
	league.PlayerCount = temp.PlayerCount

	return &league, nil
}

func (s *Store) Create(ctx context.Context, e *entities.League) (*uuid.UUID, error) {
	return nil, nil
}

func (s *Store) Update(ctx context.Context, e *entities.League) (*uuid.UUID, error) {
	return nil, nil
}

func (s *Store) Delete(ctx context.Context, e *entities.League) (*uuid.UUID, error) {
	return nil, nil
}
