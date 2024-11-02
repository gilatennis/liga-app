package court

import (
	"liga-app/entities"
	"time"
)

type VModelPageCourtDetail struct {
	User          *entities.User
	Court         *entities.Court
	AvailableDays []time.Time
}
