package league

import "liga-app/entities"

type VModelPageLeagueDetail struct {
	User                           *entities.User
	League                         *entities.League
	Team                           *entities.Team
	ShowTeamCount                  bool
	IsRatingRequired               bool
	PlayersWithRatingTooHigh       []entities.User
	PlayersWithRatingTooLow        []entities.User
	CurrentMatch                   *entities.Match
	MatchNeedsCompleteConfirmation *entities.Match
	GroupStandings                 map[string][]entities.GroupRoundRobinTeamResult
}
