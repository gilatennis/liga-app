package fixture

import (
	"liga-app/entities"
	"liga-app/infra/http/htmx/common"
)

type VModelPageAddGroup struct {
	User               *entities.User
	League             *entities.League
	Fixture            *entities.Fixture
	VModelFormAddGroup *VModelFormAddGroup
}

type VModelFormAddGroup struct {
	base common.BaseVModelForm
	Name string `form:"name" binding:"required"`
}

func (v *VModelFormAddGroup) Validate() bool {
	return v.base.Validate(*v)
}
