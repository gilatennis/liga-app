package court

import (
	"liga-app/entities"
	"time"
)

type VModelPageCourtTimes struct {
	User      *entities.User
	Court     *entities.Court
	Day       time.Time
	TimeSlots []entities.CourtTimeSlot
}
