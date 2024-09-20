package league

import "liga-app/entities"

type VModelPageLeagueList struct {
	User            *entities.User
	Leagues         []entities.League
	LeaguesNextPage int
	ShowPlayerCount bool
}
