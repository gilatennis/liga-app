package teams

import (
	"fmt"
	"liga-app/entities"
	"liga-app/enums"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func SeedLeagueTeams(db *gorm.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		users := []entities.User{}

		teamCount := c.Int("team-count")

		for i := 0; i < teamCount; i++ {
			u := entities.User{
				Email:  faker.Email(),
				Name:   faker.Name(),
				Phone:  faker.Phonenumber(),
				Gender: enums.GENDER(strings.ToUpper(faker.Gender())),
			}

			users = append(users, u)
		}

		err := db.WithContext(c.Context).Create(users).Error
		if err != nil {
			return fmt.Errorf("db.Create(users): %w", err)
		}

		return nil
	}
}

func createSinglesTeams(league *entities.League, teamCount int) []entities.Team {
	var gender enums.GENDER
	if league.GenderCategory == enums.CATEGORY_GENDER_MALE {
		gender = enums.GENDER_MALE
	}
	if league.GenderCategory == enums.CATEGORY_GENDER_FEMALE {
		gender = enums.GENDER_FEMALE
	}

	teams := []entities.Team{}
	for i := 0; i < teamCount; i++ {
		u := entities.User{
			Email:  faker.Email(),
			Name:   faker.Name(),
			Phone:  faker.Phonenumber(),
			Gender: gender,
		}

		t := entities.Team{
			Name:    u.Name,
			Captain: u,
			Players: []entities.User{u},
		}
		teams = append(teams, t)
	}

	return teams
}

func createDoublesTeams(league *entities.League, teamCount int) {
	var gender1 enums.GENDER
	var gender2 enums.GENDER
	if league.GenderCategory == enums.CATEGORY_GENDER_MALE {
		gender1 = enums.GENDER_MALE
		gender2 = enums.GENDER_MALE
	}
	if league.GenderCategory == enums.CATEGORY_GENDER_FEMALE {
		gender1 = enums.GENDER_FEMALE
		gender2 = enums.GENDER_FEMALE
	}
	if league.GenderCategory == enums.CATEGORY_GENDER_MIXED {
		gender1 = enums.GENDER_MALE
		gender2 = enums.GENDER_FEMALE
	}
}

func createBereguTeams(league *entities.League, teamCount int) {
	var gender enums.GENDER
	if league.GenderCategory == enums.CATEGORY_GENDER_MALE {
		gender = enums.GENDER_MALE
	}
	if league.GenderCategory == enums.CATEGORY_GENDER_FEMALE {
		gender = enums.GENDER_FEMALE
	}
}
