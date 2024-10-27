package court_time_slots

import (
	"errors"
	"fmt"
	"liga-app/entities"
	"liga-app/infra/gsheet"
	"time"

	"github.com/k0kubun/pp"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateCourtTimeSlots(db *gorm.DB, gsheetAPI *gsheet.API) func(c *cli.Context) error {
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

		timeslotsCourtA, err := gsheetAPI.CourtTimeSlot.GetByCourtSlug(c.Context, slugCourtA)
		if err != nil {
			return fmt.Errorf("gsheetAPI.CourtTimeSlot.GetByCourtSlug(slugCourtA, ...): %w", err)
		}

		timeslotsCourtB, err := gsheetAPI.CourtTimeSlot.GetByCourtSlug(c.Context, slugCourtB)
		if err != nil {
			return fmt.Errorf("gsheetAPI.CourtTimeSlot.GetByCourtSlug(slugCourtB, ...): %w", err)
		}

		const SLOT_DURATION time.Duration = 40 * time.Minute

		timeslots := []entities.CourtTimeSlot{}

		for _, t := range timeslotsCourtA {
			// INFO: divide time slots by 40 minutes

			totalDuration := t.End.Sub(t.Start)
			slots := int(totalDuration / SLOT_DURATION)

			start := t.Start
			for i := 0; i < slots; i++ {
				end := start.Add(SLOT_DURATION)
				ts := entities.CourtTimeSlot{
					CourtUUID: courtA.UUID,
					Start:     start,
					End:       end,
				}
				start = end

				timeslots = append(timeslots, ts)
			}
		}

		for _, t := range timeslotsCourtB {
			// INFO: divide time slots by 40 minutes

			totalDuration := t.End.Sub(t.Start)
			slots := int(totalDuration / SLOT_DURATION)

			start := t.Start
			for i := 0; i < slots; i++ {
				end := start.Add(SLOT_DURATION)
				ts := entities.CourtTimeSlot{
					CourtUUID: courtB.UUID,
					Start:     start,
					End:       end,
				}
				start = end

				timeslots = append(timeslots, ts)
			}
		}

		for _, ts := range timeslots {
			fmt.Printf("Court UUID: %s - Start: %s - End: %s \n", ts.CourtUUID, ts.Start.String(), ts.End.String())
		}

		tx := db.WithContext(c.Context).Begin()

		for _, ts := range timeslots {
			existsingTs := entities.CourtTimeSlot{}

			err := db.Where("start = ? AND end = ?", ts.Start, ts.End).First(&existsingTs).Error
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println("rec not found, should add")

				tts := ts
				err = tx.Omit(clause.Associations).Create(&tts).Error
				if err != nil {
					return fmt.Errorf("tx.Create(&tts).Error")
				}

				fmt.Println("rec created")
				continue

			} else if err != nil {
				return fmt.Errorf("db.Where(\"start = ? AND end = ?\", ts.Start, ts.End).First(&existsingTs).Error : %w", err)
			}

			pp.Println("existsingTs", existsingTs)
		}

		err = tx.Commit().Error
		if err != nil {
			return fmt.Errorf("tx.Commit().Error")
		}

		return nil
	}
}
