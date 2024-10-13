package entities

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	Base
	Number       int
	LeagueUUID   uuid.UUID  `gorm:"type:CHAR(36)"`
	League       League     `gorm:"ForeignKey:LeagueUUID;References:UUID"`
	RoundUUID    uuid.UUID  `gorm:"type:CHAR(36)"`
	Round        Round      `gorm:"ForeignKey:RoundUUID;References:UUID"`
	TeamAUUID    *uuid.UUID `gorm:"type:CHAR(36)"`
	TeamA        *Team      `gorm:"ForeignKey:TeamAUUID;References:UUID"`
	TeamBUUID    *uuid.UUID `gorm:"type:CHAR(36)"`
	TeamB        *Team      `gorm:"ForeignKey:TeamBUUID;References:UUID"`
	IsByeTeamA   bool
	IsByeTeamB   bool
	Location     string
	Time         *time.Time
	Sets         []MatchSet `gorm:"ForeignKey:MatchUUID"`
	Status       string
	FilterStatus string `gorm:"-:all"`
}
