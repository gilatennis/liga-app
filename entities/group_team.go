package entities

import "github.com/google/uuid"

type GroupTeam struct {
	Base
	GroupUUID        uuid.UUID                 `gorm:"type:CHAR(36)"`
	Group            Group                     `gorm:"ForeignKey:GroupUUID;References:UUID"`
	TeamUUID         uuid.UUID                 `gorm:"type:CHAR(36)"`
	Team             Team                      `gorm:"ForeignKey:TeamUUID;References:UUID"`
	RoundRobinResult GroupRoundRobinTeamResult `gorm:"ForeignKey:GroupUUID,TeamUUID;References:GroupUUID,TeamUUID"`
}
