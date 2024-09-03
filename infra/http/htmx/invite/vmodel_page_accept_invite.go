package invite

import "liga-app/entities"

type VModelPageAcceptInvite struct {
	User   *entities.User
	League *entities.League
	Team   *entities.Team
}
