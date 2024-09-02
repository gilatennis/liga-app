package team

import "liga-app/infra/http/htmx/common"

type VModelFormInviteViaEmail struct {
	base  common.BaseVModelForm
	Email string `validation:"email" label:"Email"`
}
