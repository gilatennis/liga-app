package entities

import "github.com/google/uuid"

type Group struct {
	Base
	FixtureUUID uuid.UUID `gorm:"type:CHAR(36)"`
	Fixture     Fixture   `gorm:"ForeignKey:FixtureUUID;References:UUID"`
	Name        string
	Rounds      []Round `gorm:"ForeignKey:GroupUUID"`
	Teams       []Team  `gorm:"many2many:group_teams;foreignKey:UUID;joinForeignKey:GroupUUID;References:UUID;joinReferences:TeamUUID"`
}
