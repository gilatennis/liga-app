package dashboard

import (
	"liga-app/infra/http/htmx/common"
	"liga-app/infra/http/htmx/common/components"
)

type VModelFormUserInfo struct {
	base common.BaseVModelForm

	ProvinceOptions []components.InputSelectOption

	Name     string `form:"name" label:"Nama" validation:"required"`
	Phone    string `form:"phone" label:"No. WA" validation:"phone"`
	Gender   string `form:"gender" label:"Gender" validation:"required"`
	Province string `form:"province" label:"Province" validation:"required"`
}

func (v *VModelFormUserInfo) Validate() bool {
	return v.base.Validate(*v)
}
