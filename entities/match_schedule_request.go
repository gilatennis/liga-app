package entities

import (
	"time"

	"github.com/google/uuid"
)

type MatchScheduleRequest struct {
	Base
	MatchUUID uuid.UUID       `gorm:"type:CHAR(36)"`
	Match     Match           `gorm:"ForeignKey:MatchUUID;References:UUID"`
	CourtUUID uuid.UUID       `gorm:"type:CHAR(36)"`
	Court     Court           `gorm:"ForeignKey:CourtUUID;References:UUID"`
	TimeSlots []CourtTimeSlot `gorm:"many2many:match_schedule_request_court_time_slots;foreignKey:UUID;joinForeignKey:MatchScheduleRequestUUID;References:UUID;joinReferences:CourtTimeSlotUUID"`

	// INFO: can use own venue and court
	VenueName string
	Address   string
	Notes     string
	Start     *time.Time
	End       *time.Time

	IsAccepted bool
}
