package services

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
)

type MATCH_FILTER_STATUS string

const (
	MATCH_FILTER_STATUS_UPCOMING  MATCH_FILTER_STATUS = "UPCOMING"
	MATCH_FILTER_STATUS_OVERDUE   MATCH_FILTER_STATUS = "OVERDUE"
	MATCH_FILTER_STATUS_COMPLETED MATCH_FILTER_STATUS = "COMPLETED"
)

type MatchFilter struct {
	Status   *MATCH_FILTER_STATUS
	UserUUID *uuid.UUID
}

type IMatchService interface {
	Get(ctx context.Context, p *common.Pagination, f *MatchFilter) ([]entities.Match, int64, error)
}
