package services

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
)

type IMatchService interface {
	Get(ctx context.Context, p *common.Pagination, f *repositories.MatchFilter) ([]entities.Match, int64, error)
}
