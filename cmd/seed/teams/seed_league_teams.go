package teams

import (
	"fmt"
	"liga-app/entities"
	"liga-app/enums"

	"github.com/go-faker/faker/v4"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedLeagueTeams(db *gorm.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		teamCount := c.Int("team-count")

		leagues := []entities.League{}

		err := db.WithContext(c.Context).Find(&leagues).Error
		if err != nil {
			return fmt.Errorf("db.WithContext(c.Context).Find(&leagues).Error")
		}

		for _, l := range leagues {
			var teams []entities.Team

			if l.ParticipationCategory == enums.CATEGORY_PARTICIPATION_SINGLES {
				teams = createSinglesTeams(&l, teamCount)
			}

			if l.ParticipationCategory == enums.CATEGORY_PARTICIPATION_DOUBLES {
				teams = createDoublesTeams(&l, teamCount)
			}

			if l.ParticipationCategory == enums.CATEGORY_PARTICIPATION_TEAMS {
				teams = createBereguTeams(&l, teamCount)
			}

			for _, t := range teams {
				for i, p := range t.Players {
					user := p
					err = db.Omit(clause.Associations).Create(&user).Error
					if err != nil {
						return err
					}
					t.Players[i] = user
				}
				t.Captain = t.Players[0]
				t.CaptainUserUUID = t.Players[0].UUID

				leagueRegistration := entities.LeagueRegistration{
					LeagueUUID:    l.UUID,
					UserUUID:      t.Captain.UUID,
					PaymentStatus: enums.PAYMENT_STATUS_PAID,
					PaymentAmount: l.RegistrationFee,
					PaymentID:     "payment-id-test",
				}
				db.Save(&leagueRegistration)

				err = db.Omit(clause.Associations).Save(&t).Error
				if err != nil {
					return err
				}

				leagueTeam := entities.LeagueTeam{
					LeagueUUID: l.UUID,
					TeamUUID:   t.UUID,
				}
				err = db.Omit(clause.Associations).Save(&leagueTeam).Error
				if err != nil {
					return err
				}

				for _, p := range t.Players {
					teamPlayer := entities.TeamPlayer{
						TeamUUID:   t.UUID,
						PlayerUUID: p.UUID,
					}
					err = db.Omit(clause.Associations).Save(&teamPlayer).Error
					if err != nil {
						return err
					}

					leaguePlayer := entities.LeaguePlayer{
						LeagueUUID: l.UUID,
						PlayerUUID: p.UUID,
					}
					err = db.Omit(clause.Associations).Save(&leaguePlayer).Error
					if err != nil {
						return err
					}
				}

			}
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

func createDoublesTeams(league *entities.League, teamCount int) []entities.Team {
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

	teams := []entities.Team{}
	for i := 0; i < teamCount; i++ {
		u1 := entities.User{
			Email:  faker.Email(),
			Name:   faker.Name(),
			Phone:  faker.Phonenumber(),
			Gender: gender1,
		}

		u2 := entities.User{
			Email:  faker.Email(),
			Name:   faker.Name(),
			Phone:  faker.Phonenumber(),
			Gender: gender2,
		}

		t := entities.Team{
			Name:    u1.Name + " & " + u2.Name,
			Captain: u1,
			Players: []entities.User{u1, u2},
		}

		teams = append(teams, t)
	}

	return teams
}

func createBereguTeams(league *entities.League, teamCount int) []entities.Team {
	var gender enums.GENDER
	if league.GenderCategory == enums.CATEGORY_GENDER_MALE {
		gender = enums.GENDER_MALE
	}
	if league.GenderCategory == enums.CATEGORY_GENDER_FEMALE {
		gender = enums.GENDER_FEMALE
	}

	teams := []entities.Team{}
	for i := 0; i < teamCount; i++ {
		teamMemberCount := 6
		teamMembers := []entities.User{}
		for j := 0; j < teamMemberCount; j++ {
			u := entities.User{
				Email:  faker.Email(),
				Name:   faker.Name(),
				Phone:  faker.Phonenumber(),
				Gender: gender,
			}
			teamMembers = append(teamMembers, u)
		}

		t := entities.Team{
			Name:    teamMembers[0].Name,
			Captain: teamMembers[0],
			Players: teamMembers,
		}

		teams = append(teams, t)
	}

	return teams
}
