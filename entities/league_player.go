package entities

import "github.com/google/uuid"

type LeaguePlayer struct {
	Base
	LeagueUUID uuid.UUID `gorm:"type:CHAR(36)"`
	League     League    `gorm:"ForeignKey:LeagueUUID;References:UUID"`
	PlayerUUID uuid.UUID `gorm:"type:CHAR(36)"`
	Player     User      `gorm:"ForeignKey:PlayerUUID;References:UUID"`
}
