package fixtures

import (
	"fmt"
	"liga-app/entities"
	"liga-app/enums"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func SeedRoundRobinFixtures(db *gorm.DB) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		leagues := []entities.League{}
		err := db.Preload("Fixtures").Find(&leagues, "season = 1").Error
		if err != nil {
			return err
		}

		for _, l := range leagues {
			if len(l.Fixtures) > 0 {
				continue
			}
			fixture := entities.Fixture{
				LeagueUUID:  l.UUID,
				Name:        "Round Robin",
				Format:      enums.FIXTURE_FORMAT_ROUND_ROBIN,
				StageNumber: 1,
			}
			err = db.Create(&fixture).Error
			if err != nil {
				return fmt.Errorf("db.Create(&fixture).Error")
			}
		}

		return nil
	}
}
