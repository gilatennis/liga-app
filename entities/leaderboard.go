package entities

import (
	"github.com/google/uuid"
)

type LeaderboardSeasonalEntry struct {
	LeagueUUID uuid.UUID
	League     League
	Rank       int
	Points     int
	Titles     int
	TeamUUID   uuid.UUID
	Team       Team
}

type LeaderboardLifetimeEntry struct {
	TeamUUID uuid.UUID
	Team     Team
	Rank     int
	Points   int
	Titles   int
}
