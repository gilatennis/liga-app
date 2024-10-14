package court_time_slot

import (
	"context"
	"fmt"
	"liga-app/entities"
	"log"
	"time"

	"google.golang.org/api/sheets/v4"
)

type API struct {
	srv *sheets.Service
}

func NewAPI(srv *sheets.Service) *API {
	return &API{srv: srv}
}

func (a *API) GetByCourtSlug(ctx context.Context, slug string) ([]entities.CourtTimeSlot, error) {
	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/16Sbt3gsI1fsKQMFCVmjf1qnVwbBfc5UDZanOITI9bqE/edit?gid=0#gid=0
	spreadsheetId := "16Sbt3gsI1fsKQMFCVmjf1qnVwbBfc5UDZanOITI9bqE"
	readRange := "dislitbang-al!A1:C30"
	resp, err := a.srv.Spreadsheets.Values.Get(spreadsheetId, readRange).MajorDimension("ROWS").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	courtTimeSlots := []entities.CourtTimeSlot{}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for i, row := range resp.Values {
			if i == 0 {
				continue
			}
			fmt.Printf("%s, %s, %s\n", row[0], row[1], row[2])
			courtSlug := row[0].(string)
			timeStartStr := row[1].(string)
			timeEndStr := row[2].(string)

			if courtSlug != slug {
				continue
			}

			timeStart, err := time.Parse(time.RFC3339, timeStartStr)
			if err != nil {
				return nil, fmt.Errorf("time.Parse(time.RFC3339, timeStartStr): %w", err)
			}

			timeEnd, err := time.Parse(time.RFC3339, timeEndStr)
			if err != nil {
				return nil, fmt.Errorf("time.Parse(time.RFC3339, timeEndStr): %w", err)
			}

			court := entities.Court{Slug: courtSlug}
			timeSlot := entities.CourtTimeSlot{
				Court: court,
				Start: timeStart,
				End:   timeEnd,
			}

			courtTimeSlots = append(courtTimeSlots, timeSlot)
		}
	}

	return courtTimeSlots, nil
}
