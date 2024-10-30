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
	GetByUUID(ctx context.Context, matchUUID uuid.UUID) (*entities.Match, error)
	RequestMatchSchedule(ctx context.Context, userUUID uuid.UUID, matchUUID uuid.UUID, timeslotUUIDs []uuid.UUID) (*uuid.UUID, error)
	RequestMatchScheduleOwnVenue(ctx context.Context, userUUID uuid.UUID, matchUUID uuid.UUID, matchScheduleRequest *entities.MatchScheduleRequest) (*uuid.UUID, error)
	ConfirmMatchSchedule(ctx context.Context, userUUID uuid.UUID, matchScheduleRequestUUID uuid.UUID) (*uuid.UUID, error)
	SaveMatchSet(ctx context.Context, userUUID uuid.UUID, matchSet *entities.MatchSet) (*uuid.UUID, error)
	CompleteMatch(ctx context.Context, userUUID, matchUUID uuid.UUID) (*uuid.UUID, error)
	ConfirmCompleteMatch(ctx context.Context, userUUID, matchUUID uuid.UUID) (*uuid.UUID, error)
	GenerateRoundRobinMatches(ctx context.Context, groupUUID uuid.UUID, isHomeAndAway bool) error
}
