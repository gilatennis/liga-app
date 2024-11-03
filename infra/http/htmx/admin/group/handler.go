package group

import (
	"liga-app/infra/http/htmx/common"
	"liga-app/interfaces/services"
)

type Handler struct {
	common.BaseHandler
	groupService services.IGroupService
}

func NewHandler(groupService services.IGroupService) *Handler {
	return &Handler{
		groupService: groupService,
	}
}
