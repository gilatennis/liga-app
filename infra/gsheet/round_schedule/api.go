package round_schedule

import (
	"fmt"
	"liga-app/utils"
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

	if utils.IsDevEnv() {
		spreadsheetId = "1vPld27g9BMhrooLTgamNhEpBklp4cAOKOzB0sBCJTZc"
	}

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
			roundNumber := row[0].(string)
			timeStartStr := row[1].(string)
			timeEndStr := row[2].(string)

			fmt.Printf("%s, %s, %s\n", roundNumber, timeStartStr, timeEndStr)

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
