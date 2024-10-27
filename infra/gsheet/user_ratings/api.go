package user_ratings

import (
	"fmt"
	"liga-app/entities"
	"strconv"

	"github.com/google/uuid"
	"google.golang.org/api/sheets/v4"
)

type API struct {
	srv *sheets.Service
}

func NewAPI(srv *sheets.Service) *API {
	return &API{srv: srv}
}

func (a *API) GetUserRatings() ([]entities.UserRating, error) {
	spreadsheetId := "1-_JaTkguErMKv9ZAlTmQFKe6pjmmxItI9bLsY90VI7E"

	readRange := "Assessments S1!A1:H"
	resp, err := a.srv.Spreadsheets.Values.Get(spreadsheetId, readRange).MajorDimension("ROWS").Do()
	if err != nil {
		return nil, fmt.Errorf("a.srv.Spreadsheets.Values.Get(spreadsheetId, readRange).MajorDimension(\"ROWS\").Do() : %w", err)
	}

	userRatings := []entities.UserRating{}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for i, row := range resp.Values {
			if i == 0 {
				continue
			}
			leagueUUIDStr := row[0].(string)
			leagueName := row[1].(string)
			teamName := row[2].(string)
			userUUIDStr := row[3].(string)
			userName := row[4].(string)
			email := row[5].(string)
			phone := row[6].(string)
			ntrpRatingStr := row[7].(string)

			fmt.Printf("%s, %s, %s, %s, %s, %s, %s, %s\n", leagueUUIDStr, leagueName, teamName, userUUIDStr, userName, email, phone, ntrpRatingStr)

			if ntrpRatingStr == "-" {
				continue
			}

			if userUUIDStr == "" {
				continue
			}

			userUUID, err := uuid.Parse(userUUIDStr)
			if err != nil {
				return nil, fmt.Errorf("uuid.Parse(userUUIDStr=%s) : %w", userUUIDStr, err)
			}

			ntrpRating, err := strconv.ParseFloat(ntrpRatingStr, 32)
			if err != nil {
				return nil, fmt.Errorf("strconv.ParseFloat(ntrpRatingStr=%s, 64) : %w", ntrpRatingStr, err)
			}

			ur := entities.UserRating{
				UserUUID:   userUUID,
				NTRPRating: float32(ntrpRating),
				Season:     1,
			}

			userRatings = append(userRatings, ur)
		}
	}

	return userRatings, nil
}
