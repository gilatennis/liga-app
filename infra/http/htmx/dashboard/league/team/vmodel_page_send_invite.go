package team

import "liga-app/entities"

type VModelPageSendInvite struct {
	User                     *entities.User
	Team                     *entities.Team
	League                   *entities.League
	InviteLink               string
	VModelFormInviteViaEmail *VModelFormInviteViaEmail
}
