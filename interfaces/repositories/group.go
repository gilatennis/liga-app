package repositories

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
)

type IGroupRepository interface {
	GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*entities.Group, error)
	Create(ctx context.Context, group *entities.Group) (*uuid.UUID, error)
	AddTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error)
	RemoveTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error)
	GetGroupTeams(ctx context.Context, groupUUID uuid.UUID) ([]entities.Team, error)
	/*
		SaveGroupRoundRobinResults(
			ctx context.Context,
			groupUUID uuid.UUID,
			groupRoundRobinResults []entities.GroupRoundRobinTeamResult,
			groupRoundRobinHeadToHeadResults []entities.GroupRoundRobinHeadToHeadResult,
		) error
	*/
}
