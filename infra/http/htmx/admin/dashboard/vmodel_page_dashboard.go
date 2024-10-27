package dashboard

import "liga-app/entities"

type VModelPageDashboard struct {
	User    *entities.User
	Leagues []entities.League
}
