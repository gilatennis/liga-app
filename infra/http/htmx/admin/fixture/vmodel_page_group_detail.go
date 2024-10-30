package fixture

import (
	"liga-app/entities"
	"liga-app/infra/http/htmx/common"
)

type VModelPageGroupDetail struct {
	User              *entities.User
	League            *entities.League
	Fixture           *entities.Fixture
	Group             *entities.Group
	TeamsAvailable    []entities.Team
	VModelFormAddTeam *VModelFormAddTeam
}

type VModelFormAddTeam struct {
	base     common.BaseVModelForm
	TeamUUID string `form:"team_uuid" binding:"required"`
}

func (v *VModelFormAddTeam) Validate() bool {
	return v.base.Validate(*v)
}
