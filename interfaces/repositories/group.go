package repositories

import (
	"context"
	"liga-app/entities"
	"liga-app/interfaces/common"

	"github.com/google/uuid"
)

type GroupFilter struct {
	FixtureUUID *uuid.UUID
	TeamUUID    *uuid.UUID
}

type IGroupRepository interface {
	Get(ctx context.Context, p *common.Pagination, f *GroupFilter) ([]entities.Group, int64, error)
	GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*entities.Group, error)
	Create(ctx context.Context, group *entities.Group) (*uuid.UUID, error)
	AddTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error)
	RemoveTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error)
	GetGroupTeams(ctx context.Context, groupUUID uuid.UUID) ([]entities.Team, error)
	SaveGroupRoundRobinResults(
		ctx context.Context,
		groupUUID uuid.UUID,
		groupRoundRobinResults []entities.GroupRoundRobinTeamResult,
		groupRoundRobinHeadToHeadResults []entities.GroupRoundRobinHeadToHeadResult,
	) error
	GetGroupRoundRobinResults(ctx context.Context, groupUUID uuid.UUID) ([]entities.GroupRoundRobinTeamResult, error)
}
