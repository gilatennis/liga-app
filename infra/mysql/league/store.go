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
		Preload("Fixtures").
		Preload("Fixtures.Groups").
		Offset(p.Offset).
		Limit(p.Limit)

	if f != nil {
		if f.UserUUID != nil {
			dbHandle = dbHandle.
				Joins("JOIN league_players ON league_players.league_uuid = leagues.uuid").
				Where("league_players.player_uuid = ?", *f.UserUUID)
		}
		if f.Season != nil {
			dbHandle = dbHandle.Where("leagues.season = ?", *f.Season)
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
		LeagueUUID                    uuid.UUID
		TeamCount                     int64
		PlayerCount                   int64
		GroupStageMatchCount          int64
		GroupStageMatchCompletedCount int64
	}
	temps := []Temp{}

	err = s.db.WithContext(ctx).
		Raw(`
      SELECT
          l.uuid AS league_uuid,
          COUNT(DISTINCT lt.uuid) AS team_count,
          COUNT(DISTINCT lp.uuid) AS player_count,
          (SELECT COUNT(DISTINCT m.uuid)
            FROM matches m
            JOIN rounds r ON m.round_uuid = r.uuid
            JOIN `+"`groups`"+` g ON r.group_uuid = g.uuid
            JOIN fixtures f ON g.fixture_uuid = f.uuid
            WHERE m.league_uuid = l.uuid AND f.format = 'ROUND_ROBIN' AND m.team_a_uuid IS NOT NULL AND m.team_b_uuid IS NOT NULL)
          AS group_stage_match_count,
          (SELECT COUNT(DISTINCT m.uuid)
            FROM matches m
            JOIN rounds r ON m.round_uuid = r.uuid
            JOIN `+"`groups`"+` g ON r.group_uuid = g.uuid
            JOIN fixtures f ON g.fixture_uuid = f.uuid
            WHERE m.league_uuid = l.uuid AND f.format = 'ROUND_ROBIN' AND m.status = 'COMPLETED'
          ) AS group_stage_match_completed_count
      FROM leagues l
      LEFT JOIN league_teams lt ON l.uuid = lt.league_uuid AND lt.deleted_at IS NULL
      LEFT JOIN league_players lp ON l.uuid = lp.league_uuid AND lp.deleted_at IS NULL
      WHERE l.uuid IN ?
      GROUP BY l.uuid;
    `, leagueUUIDs).
		Scan(&temps).
		Error
	if err != nil {
		return nil, 0, err
	}

	for i, l := range leagues {
		l.TeamCount = temps[i].TeamCount
		l.PlayerCount = temps[i].PlayerCount
		l.GroupStageMatchCount = temps[i].GroupStageMatchCount
		l.GroupStageMatchCompletedCount = temps[i].GroupStageMatchCompletedCount
		leagues[i] = l
	}

	return leagues, count, nil
}

func (s *Store) GetByUUID(ctx context.Context, leagueUUID uuid.UUID) (*entities.League, error) {
	league := entities.League{}
	err := s.db.Where("uuid = ?", leagueUUID).
		Preload("Wilayah").
		Preload("Fixtures").
		Preload("Fixtures.Groups").
		First(&league).
		Error
	if err != nil {
		return nil, err
	}

	type Temp struct {
		LeagueUUID                    uuid.UUID
		TeamCount                     int64
		PlayerCount                   int64
		GroupStageMatchCount          int64
		GroupStageMatchCompletedCount int64
	}
	temp := Temp{}
	err = s.db.WithContext(ctx).
		Raw(`
      select
        l.uuid AS league_uuid,
        COUNT(DISTINCT lt.uuid) AS team_count,
        COUNT(DISTINCT lp.uuid) AS player_count,
        (SELECT COUNT(DISTINCT m.uuid)
          FROM matches m
          JOIN rounds r ON m.round_uuid = r.uuid
          JOIN `+"`groups`"+` g ON r.group_uuid = g.uuid
          JOIN fixtures f ON g.fixture_uuid = f.uuid
          WHERE m.league_uuid = l.uuid AND f.format = 'ROUND_ROBIN' AND m.team_a_uuid IS NOT NULL AND m.team_b_uuid IS NOT NULL)
        AS group_stage_match_count,
        (SELECT COUNT(DISTINCT m.uuid)
          FROM matches m
          JOIN rounds r ON m.round_uuid = r.uuid
          JOIN `+"`groups`"+` g ON r.group_uuid = g.uuid
          JOIN fixtures f ON g.fixture_uuid = f.uuid
          WHERE m.league_uuid = l.uuid AND f.format = 'ROUND_ROBIN' AND m.status = 'COMPLETED'
        ) AS group_stage_match_completed_count
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
	league.GroupStageMatchCount = temp.GroupStageMatchCount
	league.GroupStageMatchCompletedCount = temp.GroupStageMatchCompletedCount

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
