package entities

import (
	"time"

	"github.com/google/uuid"
)

type CourtTimeSlot struct {
	Base
	CourtUUID uuid.UUID `gorm:"type:CHAR(36)"`
	Court     Court     `gorm:"ForeignKey:CourtUUID;References:UUID"`
	Start     time.Time
	End       time.Time
	MatchUUID *uuid.UUID `gorm:"type:CHAR(36)"`
	Match     *Match     `gorm:"ForeignKey:MatchUUID;References:UUID"`
}
