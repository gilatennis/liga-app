package match

import "liga-app/entities"

type VModelPageMatchList struct {
	User            *entities.User
	Matches         []entities.Match
	MatchesNextPage int
}
