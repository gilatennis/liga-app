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
	End       time.Time `form:"end"`
}

func (v *VModelFormMatchScheduleOwnVenue) Validate() bool {
	isValid := v.base.Validate(*v)

	startDate := v.Start.Day()
	endDate := v.End.Day()

	if startDate != endDate {
		v.base.FormErrors["end"] = "waktu mulai dan waktu selesai harus dalam hari yang sama"
		isValid = false
	}

	return isValid
}
