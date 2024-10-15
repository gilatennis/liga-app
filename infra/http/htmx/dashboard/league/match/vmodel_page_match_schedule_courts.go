package match

import (
	"liga-app/entities"
)

type VModelPageMatchScheduleCourts struct {
	User   *entities.User
	League *entities.League
	Match  *entities.Match
	Venue  *entities.Venue
}
