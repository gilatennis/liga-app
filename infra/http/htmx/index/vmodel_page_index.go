package index

import "liga-app/entities"

type VModelPageIndex struct {
	RuleBook       *entities.RuleBook
	User           *entities.User
	Leagues        []entities.League
	GroupStandings map[string][]entities.GroupRoundRobinTeamResult
}
