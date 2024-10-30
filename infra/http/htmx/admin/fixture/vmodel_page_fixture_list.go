package fixture

import "liga-app/entities"

type VModelPageFixtureList struct {
	User     *entities.User
	League   *entities.League
	Fixtures []entities.Fixture
}
