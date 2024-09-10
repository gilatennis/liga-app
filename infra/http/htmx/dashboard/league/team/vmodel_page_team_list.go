package team

import "liga-app/entities"

type VModelPageTeamList struct {
	User          *entities.User
	League        *entities.League
	Teams         []entities.Team
	TeamsNextPage int
}
