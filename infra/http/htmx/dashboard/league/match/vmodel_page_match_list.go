package match

import "liga-app/entities"

type VModelPageMatchList struct {
	User            *entities.User
	League          *entities.League
	Matches         []entities.Match
	MatchesNextPage int
}
