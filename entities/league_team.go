package entities

import "github.com/google/uuid"

type LeagueTeam struct {
	Base
	LeagueUUID uuid.UUID `gorm:"type:CHAR(36)"`
	League     League    `gorm:"ForeignKey:LeagueUUID;References:UUID"`
	TeamUUID   uuid.UUID `gorm:"type:CHAR(36)"`
	Team       Team      `gorm:"ForeignKey:TeamUUID;References:UUID"`
}
