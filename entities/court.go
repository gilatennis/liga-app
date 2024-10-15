package entities

import "github.com/google/uuid"

type Court struct {
	Base
	VenueUUID uuid.UUID `gorm:"type:CHAR(36)"`
	Venue     Venue     `gorm:"ForeignKey:VenueUUID;References:UUID"`
	Name      string
	Slug      string
	TimeSlots []CourtTimeSlot `gorm:"ForeignKey:CourtUUID"`
	Image     string
}
