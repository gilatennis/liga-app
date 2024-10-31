package fixture

import "liga-app/entities"

type VModelPageGroupMatches struct {
	User            *entities.User
	League          *entities.League
	Fixture         *entities.Fixture
	Group           *entities.Group
	Matches         []entities.Match
	MatchesNextPage int
}
