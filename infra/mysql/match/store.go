package match

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
	"time"

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

func (s *Store) Get(ctx context.Context, p *common.Pagination, f *repositories.MatchFilter) ([]entities.Match, int64, error) {
	var matches []entities.Match
	var count int64

	queryDB := s.db.WithContext(ctx).
		Model(&entities.Match{}).
		Preload("Sets").
		Preload("Round.Group.Fixture").
		Joins("JOIN rounds on matches.round_uuid = rounds.uuid").
		Joins("JOIN `groups` on rounds.group_uuid = `groups`.uuid").
		Joins("JOIN fixtures on `groups`.fixture_uuid = fixtures.uuid").
		Joins("JOIN leagues on fixtures.league_uuid = leagues.uuid")

	for _, o := range p.OrderBy {
		queryDB = queryDB.Order(fmt.Sprintf("%s %s", o.Column, o.Direction))
	}

	if f != nil {
		if f.RoundIsAfter != nil {
			queryDB = queryDB.Where("rounds.start >= ?", f.RoundIsAfter)
		}
		if f.RoundIsBefore != nil {
			queryDB = queryDB.Where("rounds.start <= ?", f.RoundIsBefore)
		}
		if f.UserUUID != nil {
			queryDB = queryDB.
				Joins("LEFT JOIN team_players tpa on matches.team_a_uuid = tpa.team_uuid AND tpa.player_uuid = ?", f.UserUUID).
				Joins("LEFT JOIN team_players tpb on matches.team_b_uuid = tpb.team_uuid AND tpb.player_uuid = ?", f.UserUUID).
				Where("tpa.player_uuid = ? OR tpb.player_uuid = ?", f.UserUUID, f.UserUUID)
		}

		if f.Status != nil {
			queryDB = queryDB.Where("status = ?", *f.Status)
		}

		if f.GroupUUID != nil {
			queryDB = queryDB.Where("`groups`.uuid = ?", *f.GroupUUID)
		}

		if f.LeagueUUID != nil {
			queryDB = queryDB.Where("leagues.uuid = ?", *f.LeagueUUID)
		}

		if f.IsScheduleConfirmed != nil {
			queryDB = queryDB.
				Joins("LEFT JOIN match_schedule_requests on matches.uuid = match_schedule_requests.match_uuid").
				Where("match_schedule_requests.is_accepted = ?", *f.IsScheduleConfirmed)
		}

		if f.IsComplete != nil {
			queryDB = queryDB.Where("is_complete = ?", *f.IsComplete)
		}

		if f.IsCompleteConfirmed != nil {
			queryDB = queryDB.Where("is_complete_confirmed = ?", *f.IsCompleteConfirmed)
		}
	}

	err := queryDB.Debug().
		Preload("League").Preload("TeamA").Preload("TeamB").Preload("Round").Preload("MatchScheduleRequest").
		Offset(p.Offset).
		Limit(p.Limit).
		Find(&matches).
		Error
	if err != nil {
		return nil, 0, fmt.Errorf("s.db.WithContext(ctx).Find(...): %w", err)
	}

	err = queryDB.Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("s.db.WithContext(ctx).Model(&entities.Match{}).Count(&count).Error: %w", err)
	}

	fmt.Println("store : offset : ", p.Offset)
	fmt.Println("store : limit : ", p.Limit)
	fmt.Println("store : count : ", count)

	return matches, count, nil
}

func (s *Store) GetByUUID(ctx context.Context, matchUUID uuid.UUID) (*entities.Match, error) {
	match := entities.Match{}

	err := s.db.Where("uuid = ?", matchUUID).
		Preload("League").
		Preload("TeamA").
		Preload("TeamA.Captain").
		Preload("TeamB").
		Preload("TeamB.Captain").
		Preload("Round").
		Preload("MatchScheduleRequest").
		Preload("MatchScheduleRequest.Court").
		Preload("MatchScheduleRequest.Court.Venue").
		Preload("MatchScheduleRequest.TimeSlots").
		Preload("Sets").
		First(&match).Error
	if err != nil {
		return nil, fmt.Errorf("s.db.Where(\"uuid = ?\", matchUUID=%s).First(&match).Error : %w", matchUUID, err)
	}

	return &match, nil
}

func (s *Store) Update(ctx context.Context, match *entities.Match) (*uuid.UUID, error) {
	match.UpdatedAt = time.Now()
	err := s.db.WithContext(ctx).Omit(clause.Associations).Save(match).Error
	if err != nil {
		return nil, fmt.Errorf("s.db.Save(match).Error : %w", err)
	}

	return &match.UUID, nil
}

func (s *Store) CreateGeneratedMatches(ctx context.Context, rounds []entities.Round, leagueUUID, groupUUID uuid.UUID) error {
	tx := s.db.WithContext(ctx).Begin()

	for _, r := range rounds {
		r.GroupUUID = groupUUID

		err := tx.Omit(clause.Associations).Create(&r).Error
		if err != nil {
			return fmt.Errorf("tx.Omit(clause.Associations).Create(r).Error: %w", err)
		}

		for j, m := range r.Matches {
			m.LeagueUUID = leagueUUID
			m.RoundUUID = r.UUID

			r.Matches[j] = m
		}

		err = tx.Omit(clause.Associations).Create(r.Matches).Error
		if err != nil {
			return fmt.Errorf("tx.Omit(clause.Associations).Create(r.Matches).Error : %w", err)
		}
	}

	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("tx.Commit(): %w", err)
	}

	return nil
}

func (s *Store) GetPlayerMatchStats(ctx context.Context, userUUID uuid.UUID) (*entities.MatchStats, error) {
	matchStats := entities.MatchStats{}

	err := s.db.WithContext(ctx).
		Raw(`
      SELECT
        user_uuid,
        COUNT(*) as played,
        SUM(CASE when sets_win > sets_lose then 1 else 0 END) as win,
        SUM(CASE when sets_lose > sets_win then 1 else 0 END) as lose
      FROM
      (SELECT 
        -- m.*,
        m.uuid,
        -- m.team_a_uuid,
        ta.name as team_a_name,
        -- m.team_b_uuid,
        tb.name as team_b_name,
        -- (case when m.team_a_uuid = tp.team_uuid then m.team_a_uuid else m.team_b_uuid END) as team_uuid,
        t.name as team_name,
        u.uuid as user_uuid,
        u.name as user_name,
        SUM(ms.team_a_score),
        SUM(ms.team_b_score),
        SUM(case 
        when (m.team_a_uuid = tp.team_uuid AND ms.team_a_score > ms.team_b_score) then 1 
          when (m.team_a_uuid = tp.team_uuid AND ms.team_b_score > ms.team_a_score) then 0
          when (m.team_b_uuid = tp.team_uuid AND ms.team_a_score > ms.team_b_score) then 0 
          else 1 END
        ) AS sets_win,
        SUM(case 
        when (m.team_a_uuid = tp.team_uuid AND ms.team_a_score > ms.team_b_score) then 0 
          when (m.team_a_uuid = tp.team_uuid AND ms.team_b_score > ms.team_a_score) then 1
          when (m.team_b_uuid = tp.team_uuid AND ms.team_a_score > ms.team_b_score) then 1 
          else 0 END
        ) AS sets_lose
      FROM 
        users u 
          JOIN team_players tp on tp.player_uuid = u.uuid 
          JOIN teams t on tp.team_uuid = t.uuid
          JOIN matches m on tp.team_uuid = (case when (m.team_a_uuid = tp.team_uuid) then m.team_a_uuid else m.team_b_uuid END)
          JOIN match_sets ms on ms.match_uuid = m.uuid
          JOIN teams ta on m.team_a_uuid = ta.uuid
          JOIN teams tb on m.team_b_uuid = tb.uuid
      WHERE u.uuid = ?
      GROUP BY m.uuid, t.name, u.uuid) inner_table
      GROUP BY user_uuid;
    `, userUUID).
		Scan(&matchStats).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(ctx).Raw(...).Error : %w", err)
	}

	return &matchStats, nil
}
