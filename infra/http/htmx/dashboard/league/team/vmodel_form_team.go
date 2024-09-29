package team

import (
	"liga-app/infra/http/htmx/common"
	"mime/multipart"

	"github.com/google/uuid"
)

type VModelFormTeam struct {
	base                  common.BaseVModelForm
	LeagueUUID            uuid.UUID
	TeamUUID              uuid.UUID
	ProfilePicture        *multipart.FileHeader `form:"profile_picture" label:"Profile Picture"`
	ProfilePictureCurrent string
	Name                  string `form:"name" label:"Nama Tim" validation:"required"`
}

func (v *VModelFormTeam) Validate() bool {
	return v.base.Validate(*v)
}
