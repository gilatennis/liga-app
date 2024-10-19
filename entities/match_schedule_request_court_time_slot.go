package entities

import "github.com/google/uuid"

type MatchScheduleRequestCourtTimeSlot struct {
	MatchScheduleRequestUUID uuid.UUID `gorm:"type:CHAR(36)"`
	CourtTimeSlotUUID        uuid.UUID `gorm:"type:CHAR(36)"`
}
