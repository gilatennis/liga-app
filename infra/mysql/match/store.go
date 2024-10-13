package match

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"

	"gorm.io/gorm"
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
		Joins("JOIN rounds on matches.round_uuid = rounds.uuid")

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
