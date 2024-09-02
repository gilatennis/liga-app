package entities

import "github.com/google/uuid"

type TeamPlayer struct {
	Base
	TeamUUID   uuid.UUID `gorm:"type:CHAR(36)"`
	Team       Team      `gorm:"ForeignKey:TeamUUID;References:UUID"`
	PlayerUUID uuid.UUID `gorm:"type:CHAR(36)"`
	Player     User      `gorm:"ForeignKey:PlayerUUID;References:UUID"`
}
