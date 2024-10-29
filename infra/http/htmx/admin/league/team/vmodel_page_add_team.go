package team

import (
	"liga-app/entities"
	"liga-app/infra/http/htmx/common"
)

type VModelPageAddTeam struct {
	User              *entities.User
	League            *entities.League
	Players           []entities.User
	VModelFormAddTeam *VModelFormAddTeam
}

type VModelFormAddTeam struct {
	base          common.BaseVModelForm
	TeamName      string   `form:"team_name"`
	PlayerUUIDs   []string `form:"[]player_uuids"`
	PaymentAmount float64  `form:"payment_amount"`
	PaymentID     string   `form:"payment_id"`
}
