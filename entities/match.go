package entities

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	Base
	LeagueUUID uuid.UUID `gorm:"type:CHAR(36)"`
	League     League    `gorm:"ForeignKey:LeagueUUID;References:UUID"`
	TeamAUUID  uuid.UUID `gorm:"type:CHAR(36)"`
	TeamA      Team      `gorm:"ForeignKey:TeamAUUID;References:UUID"`
	TeamBUUID  uuid.UUID `gorm:"type:CHAR(36)"`
	TeamB      Team      `gorm:"ForeignKey:TeamBUUID;References:UUID"`
	Location   string
	Time       *time.Time
	Sets       []MatchSet `gorm:"ForeignKey:MatchUUID"`
	Status     string
}
