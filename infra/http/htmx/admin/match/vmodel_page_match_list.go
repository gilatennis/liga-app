package match

import (
	"liga-app/entities"

	"github.com/google/uuid"
)

type VModelPageMatchList struct {
	User            *entities.User
	LeagueUUID      uuid.UUID
	FixtureUUID     uuid.UUID
	GroupUUID       uuid.UUID
	Matches         []entities.Match
	MatchesNextPage int
}
