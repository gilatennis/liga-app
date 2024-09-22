package profile

import (
	"liga-app/infra/http/htmx/common"
	"liga-app/infra/http/htmx/common/components"
	"mime/multipart"
)

type VModelFormProfile struct {
	base common.BaseVModelForm

	ProvinceOptions []components.InputSelectOption

	ProfilePicture        *multipart.FileHeader `form:"profile_picture" label:"Profile Picture"`
	ProfilePictureCurrent string

	Name     string `form:"name" label:"Nama" validation:"required"`
	Phone    string `form:"phone" label:"No. WA" validation:"phone"`
	Gender   string `form:"gender" label:"Gender" validation:"required"`
	Province string `form:"province" label:"Province" validation:"required"`
}

func (v *VModelFormProfile) Validate() bool {
	return v.base.Validate(*v)
}
