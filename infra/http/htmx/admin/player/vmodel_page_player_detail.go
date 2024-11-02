package player

import (
	"liga-app/entities"
	"liga-app/infra/http/htmx/common"
	"liga-app/infra/http/htmx/common/components"
	"mime/multipart"
)

type VModelPagePlayerDetail struct {
	User             *entities.User
	Player           *entities.User
	VModelFormPlayer *VModelFormPlayer
}

type VModelFormPlayer struct {
	base common.BaseVModelForm

	ProvinceOptions []components.InputSelectOption

	ProfilePicture        *multipart.FileHeader `form:"profile_picture" label:"Profile Picture"`
	ProfilePictureCurrent string

	Name  string `form:"name" label:"Nama" validation:"required"`
	Phone string `form:"phone" label:"No. WA" validation:"phone"`
	// Gender        string   `form:"gender" label:"Gender" validation:"required"`
	Province      string   `form:"province" label:"Province" validation:"required"`
	RatingSeason1 *float32 `form:"rating_season_1" label:"Rating Season 1"`
}

func (v *VModelFormPlayer) Validate() bool {
	return v.base.Validate(*v)
}
