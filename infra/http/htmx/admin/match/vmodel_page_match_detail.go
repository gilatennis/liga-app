package match

import "liga-app/entities"

type VModelPageMatchDetail struct {
	User   *entities.User
	League *entities.League
	Match  *entities.Match
}
