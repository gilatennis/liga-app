package entities

import (
	"liga-app/enums"

	"github.com/google/uuid"
)

type GroupRoundRobinTeamResult struct {
	Rank        int
	GroupUUID   uuid.UUID
	TeamUUID    uuid.UUID
	Players     []User
	Score       int
	MatchPlayed int
	MatchWon    int
	MatchLost   int
}

type GroupRoundRobinHeadToHeadResult struct {
	GroupUUID        uuid.UUID
	TeamUUID         uuid.UUID
	TeamOpponentUUID uuid.UUID
	Result           enums.MatchResult
}
