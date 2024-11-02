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
	}

	err := queryDB.Debug().
		Preload("League").Preload("TeamA").Preload("TeamB").Preload("Round").
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
