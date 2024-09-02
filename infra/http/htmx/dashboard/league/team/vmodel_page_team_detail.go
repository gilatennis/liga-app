package team

import "liga-app/entities"

type VModelPageTeamDetail struct {
	User           *entities.User
	League         *entities.League
	Team           *entities.Team
	VModelFormTeam *VModelFormTeam
}
