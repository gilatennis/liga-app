package fixtures

import (
	"fmt"
	"liga-app/entities"
	"liga-app/enums"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func SeedKnockoutFixtures(db *gorm.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		leagues := []entities.League{}
		err := db.Preload("Fixtures").Find(&leagues, "season = 1").Error
		if err != nil {
			return err
		}

		for _, l := range leagues {
			for _, f := range l.Fixtures {
				if f.Format == enums.FIXTURE_FORMAT_KNOCKOUT {
					continue
				}
			}
			type Temp struct {
				LeagueUUID  uuid.UUID
				TeamCount   int64
				PlayerCount int64
			}
			temp := Temp{}
			err = db.
				Raw(`
          select
            l.uuid AS league_uuid,
            COUNT(DISTINCT lt.uuid) AS team_count,
            COUNT(DISTINCT lp.uuid) AS player_count
          from leagues l
            left join league_teams lt on l.uuid = lt.league_uuid 
            left join league_players lp on l.uuid = lp.league_uuid
          WHERE l.uuid = ? AND lt.deleted_at IS NULL AND lp.deleted_at IS NULL
          GROUP BY l.uuid
        `, l.UUID).
				Scan(&temp).
				Error
			if err != nil {
				return err
			}

			tx := db.Begin()

			fixture := entities.Fixture{
				LeagueUUID:  l.UUID,
				Name:        "Knockout",
				Format:      enums.FIXTURE_FORMAT_KNOCKOUT,
				StageNumber: 2,
			}
			err = tx.Create(&fixture).Error
			if err != nil {
				return fmt.Errorf("db.Create(&fixture).Error")
			}

			group := entities.Group{
				FixtureUUID: fixture.UUID,
				Name:        "Default",
			}
			err = tx.Create(&group).Error
			if err != nil {
				return fmt.Errorf("db.Create(&group).Error")
			}

			rounds := []entities.Round{}
			roundCount := 0
			if temp.TeamCount <= 4 {
				roundCount = 1
			} else {
				roundCount = 2
			}
			for i := 0; i < roundCount; i++ {
				r := entities.Round{
					GroupUUID: group.UUID,
					Number:    i + 1,
				}
				rounds = append(rounds, r)
			}
			err = tx.Create(&rounds).Error
			if err != nil {
				return fmt.Errorf("db.Create(&rounds).Error")
			}

			err = tx.Commit().Error
			if err != nil {
				return fmt.Errorf("tx.Commit().Error")
			}
		}

		return nil
	}
}
