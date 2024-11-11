package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/interfaces/common"
	"time"

	"github.com/google/uuid"
)

type MatchFilter struct {
	RoundIsAfter        *time.Time
	RoundIsBefore       *time.Time
	UserUUID            *uuid.UUID
	Status              *enums.MATCH_STATUS
	GroupUUID           *uuid.UUID
	LeagueUUID          *uuid.UUID
	IsScheduleConfirmed *bool
	IsComplete          *bool
	IsCompleteConfirmed *bool
}

type IMatchRepository interface {
	Get(ctx context.Context, p *common.Pagination, f *MatchFilter) ([]entities.Match, int64, error)
	GetByUUID(ctx context.Context, matchUUID uuid.UUID) (*entities.Match, error)
	Update(ctx context.Context, m *entities.Match) (*uuid.UUID, error)
	CreateGeneratedMatches(ctx context.Context, rounds []entities.Round, leagueUUID, groupUUID uuid.UUID) error
	GetPlayerMatchStats(ctx context.Context, userUUID uuid.UUID) (*entities.MatchStats, error)
}
