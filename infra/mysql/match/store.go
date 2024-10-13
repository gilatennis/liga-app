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
