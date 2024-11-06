package entities

import (
	"liga-app/enums"

	"github.com/google/uuid"
)

type GroupRoundRobinTeamResult struct {
	Base
	GroupUUID   uuid.UUID `gorm:"type:CHAR(36)"`
	Group       Group     `gorm:"ForeignKey:GroupUUID;References:UUID"`
	TeamUUID    uuid.UUID `gorm:"type:CHAR(36)"`
	Team        Team      `gorm:"ForeignKey:TeamUUID;References:UUID"`
	Rank        int
	Score       int
	MatchPlayed int
	MatchWon    int
	MatchLost   int
	GamesWon    int
	GamesLost   int
}

type GroupRoundRobinHeadToHeadResult struct {
	Base
	GroupUUID        uuid.UUID `gorm:"type:CHAR(36)"`
	Group            Group     `gorm:"ForeignKey:GroupUUID;References:UUID"`
	TeamUUID         uuid.UUID `gorm:"type:CHAR(36)"`
	Team             Team      `gorm:"ForeignKey:TeamUUID;References:UUID"`
	TeamOpponentUUID uuid.UUID
	TeamOpponent     Team `gorm:"ForeignKey:TeamOpponentUUID;References:UUID"`
	Result           enums.MatchResult
}
