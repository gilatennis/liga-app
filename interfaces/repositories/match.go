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
	RoundIsAfter  *time.Time
	RoundIsBefore *time.Time
	UserUUID      *uuid.UUID
	Status        *enums.MATCH_STATUS
}

type IMatchRepository interface {
	Get(ctx context.Context, p *common.Pagination, f *MatchFilter) ([]entities.Match, int64, error)
}
