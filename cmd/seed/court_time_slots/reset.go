package court_time_slots

import (
	"fmt"
	"liga-app/entities"
	"liga-app/infra/gsheet"
	"time"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ResetCourtTimeSlots(db *gorm.DB, gsheetAPI *gsheet.API) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		slugCourtA := "dislitbang-al-a"
		slugCourtB := "dislitbang-al-b"

		var courtA entities.Court
		var courtB entities.Court

		err := db.WithContext(c.Context).Where("slug = ?", slugCourtA).First(&courtA).Error
		if err != nil {
			return fmt.Errorf("db.Where(\"slug = 'dislitbang-al-a'\").First(&court).Error")
		}

		err = db.WithContext(c.Context).Where("slug = ?", slugCourtB).First(&courtB).Error
		if err != nil {
			return fmt.Errorf("db.Where(\"slug = 'dislitbang-al-a'\").First(&court).Error")
		}

		tx := db.WithContext(c.Context).Begin()

		err = tx.Delete(&entities.CourtTimeSlot{}, "court_uuid = ?", courtA.UUID).Error
		if err != nil {
			return fmt.Errorf("tx.Delete(&entities.CourtTimeSlot{}, \"court_uuid = ?\", courtA.UUID).Error : %w", err)
		}

		err = tx.Delete(&entities.CourtTimeSlot{}, "court_uuid = ?", courtB.UUID).Error
		if err != nil {
			return fmt.Errorf("tx.Delete(&entities.CourtTimeSlot{}, \"court_uuid = ?\", courtB.UUID).Error : %w", err)
		}

		rangeStart := "2024-11-01T00:00:00+07:00"
		rangeEnd := "2024-12-30T23:59:59+07:00"
		timeRangeStart, _ := time.Parse(time.RFC3339, rangeStart)
		timeRangeEnd, _ := time.Parse(time.RFC3339, rangeEnd)

		timeslotsCourtA, err := gsheetAPI.CourtTimeSlot.GetByCourtSlug(c.Context, slugCourtA, timeRangeStart, timeRangeEnd)
		if err != nil {
			return fmt.Errorf("gsheetAPI.CourtTimeSlot.GetByCourtSlug(slugCourtA, ...): %w", err)
		}

		for i, t := range timeslotsCourtA {
			t.CourtUUID = courtA.UUID
			timeslotsCourtA[i] = t
		}

		timeslotsCourtB, err := gsheetAPI.CourtTimeSlot.GetByCourtSlug(c.Context, slugCourtB, timeRangeStart, timeRangeEnd)
		if err != nil {
			return fmt.Errorf("gsheetAPI.CourtTimeSlot.GetByCourtSlug(slugCourtB, ...): %w", err)
		}

		for i, t := range timeslotsCourtB {
			t.CourtUUID = courtA.UUID
			timeslotsCourtB[i] = t
		}

		timeslots := append(timeslotsCourtA, timeslotsCourtB...)

		err = tx.Omit(clause.Associations).Create(timeslots).Error
		if err != nil {
			return fmt.Errorf("tx.Create(timeslots).Error")
		}

		err = tx.Commit().Error
		if err != nil {
			return fmt.Errorf("tx.Commit().Error")
		}

		return nil
	}
}
