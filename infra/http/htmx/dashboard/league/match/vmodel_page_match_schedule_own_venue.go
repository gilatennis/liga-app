package match

import (
	"liga-app/entities"
	"liga-app/infra/http/htmx/common"
	"time"
)

type VModelPageMatchScheduleOwnVenue struct {
	User                            *entities.User
	League                          *entities.League
	Match                           *entities.Match
	VModelFormMatchScheduleOwnVenue *VModelFormMatchScheduleOwnVenue
}

type VModelFormMatchScheduleOwnVenue struct {
	base      common.BaseVModelForm
	VenueName string    `form:"venue_name"`
	Address   string    `form:"address"`
	Notes     string    `form:"notes"`
	Start     time.Time `form:"start"`
}

func (v *VModelFormMatchScheduleOwnVenue) Validate() bool {
	isValid := v.base.Validate(*v)

	return isValid
}
