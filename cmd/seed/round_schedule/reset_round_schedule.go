package round_schedule

import (
	"fmt"
	"liga-app/entities"
	"liga-app/infra/gsheet"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func ResetRoundSchedules(db *gorm.DB, gsheetAPI *gsheet.API) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		roundSchedules, err := gsheetAPI.RoundSchedule.GetRoundSchedule()
		if err != nil {
			return fmt.Errorf("gsheetAPI.RoundSchedule.GetRoundSchedule() : %w", err)
		}

		for i, schedule := range roundSchedules {
			start := schedule[0]
			end := schedule[1]
			roundNumber := i + 1

			err := db.WithContext(c.Context).
				Model(&entities.Round{}).
				Where("number = ?", roundNumber).
				UpdateColumn("start", start).
				UpdateColumn("end", end).
				Error
			if err != nil {
				return fmt.Errorf("db.WithContext()...UpdateColumn(...)... : %w", err)
			}
		}

		return nil
	}
}
