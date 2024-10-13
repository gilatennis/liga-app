package round_schedule

import (
	"fmt"
	"time"

	"google.golang.org/api/sheets/v4"
)

type API struct {
	srv *sheets.Service
}

func NewAPI(srv *sheets.Service) *API {
	return &API{srv: srv}
}

func (a *API) GetRoundSchedule() ([][2]time.Time, error) {
	spreadsheetId := "18GD0bdFLf1wj_b-KMBQ-MjAsTp_8zcpkKeeTVz2oBjI"

	readRange := "Round Schedule!A1:E"
	resp, err := a.srv.Spreadsheets.Values.Get(spreadsheetId, readRange).MajorDimension("ROWS").Do()
	if err != nil {
		return nil, fmt.Errorf("a.srv.Spreadsheets.Values.Get(spreadsheetId, readRange).MajorDimension(\"ROWS\").Do() : %w", err)
	}

	roundSchedules := [][2]time.Time{}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for i, row := range resp.Values {
			if i == 0 {
				continue
			}
			// fmt.Printf("%s, %s, %s\n", row[0], row[1], row[2])
			roundNumber := row[0].(string)
			timeStartStr := row[1].(string)
			timeEndStr := row[2].(string)

			_ = roundNumber

			timeStart, err := time.Parse(time.RFC3339, timeStartStr)
			if err != nil {
				return nil, fmt.Errorf("time.Parse(time.RFC3339, timeStartStr): %w", err)
			}

			timeEnd, err := time.Parse(time.RFC3339, timeEndStr)
			if err != nil {
				return nil, fmt.Errorf("time.Parse(time.RFC3339, timeEndStr): %w", err)
			}

			schedule := [2]time.Time{timeStart, timeEnd}
			roundSchedules = append(roundSchedules, schedule)
		}
	}

	return roundSchedules, nil
}
