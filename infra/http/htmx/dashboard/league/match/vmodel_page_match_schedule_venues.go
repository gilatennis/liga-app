package match

import "liga-app/entities"

type VModelPageMatchScheduleVenues struct {
	User   *entities.User
	League *entities.League
	Match  *entities.Match
	Venues []entities.Venue
}
