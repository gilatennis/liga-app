package match

import (
	"liga-app/entities"
	"time"
)

type VModelPageMatchScheduleTimes struct {
	User      *entities.User
	League    *entities.League
	Match     *entities.Match
	Venue     *entities.Venue
	Court     *entities.Court
	Day       time.Time
	TimeSlots []entities.CourtTimeSlot
}
