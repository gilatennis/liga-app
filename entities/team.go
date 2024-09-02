package entities

import "github.com/google/uuid"

type Team struct {
	Base
	Name            string
	CaptainUserUUID uuid.UUID `gorm:"type:CHAR(36)"`
	Captain         User      `gorm:"ForeignKey:CaptainUserUUID;References:UUID"`
	Players         []User    `gorm:"many2many:team_players;foreignKey:UUID;joinForeignKey:TeamUUID;References:UUID;joinReferences:PlayerUUID"`
}
