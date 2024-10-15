package match

import (
	"liga-app/entities"
	"time"
)

type VModelPageMatchScheduleDays struct {
	User          *entities.User
	League        *entities.League
	Match         *entities.Match
	Venue         *entities.Venue
	Court         *entities.Court
	AvailableDays []time.Time
}
