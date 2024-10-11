package entities

import (
	"time"

	"github.com/google/uuid"
)

type MatchScheduleRequest struct {
	Base
	CourtUUID uuid.UUID `gorm:"type:CHAR(36)"`
	Court     Court     `gorm:"ForeignKey:CourtUUID;References:UUID"`
	Start     time.Time
	End       time.Time
}
