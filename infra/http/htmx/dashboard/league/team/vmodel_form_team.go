package team

import "liga-app/infra/http/htmx/common"

type VModelFormTeam struct {
	base common.BaseVModelForm
	Name string `validation:"required" label:"Nama Tim"`
}
