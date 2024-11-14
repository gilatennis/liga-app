package match

import "liga-app/entities"

type VModelPageMatchList struct {
	User            *entities.User
	League          *entities.League
	Teams           []entities.Team
	Matches         []entities.Match
	MatchesNextPage int
}
