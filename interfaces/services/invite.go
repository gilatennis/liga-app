package services

import (
	"context"

	"github.com/google/uuid"
)

type IInviteService interface {
	AcceptLeagueTeamInvite(ctx context.Context, userUUID, leagueUUID, teamUUID uuid.UUID) error
}
