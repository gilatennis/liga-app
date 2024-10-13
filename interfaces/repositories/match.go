package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
)

type MatchFilter struct {
	TeamUUID *uuid.UUID
}

type IMatchRepository interface {
	Get(ctx context.Context, p *common.Pagination, f *MatchFilter) ([]entities.Match, int64, error)
}
