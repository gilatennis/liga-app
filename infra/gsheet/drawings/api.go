package drawings

// 1OT-r1U4_xeRmg53Wj62ujbfto0aaURXynyu80ka2vOQ

import (
	"fmt"
	"liga-app/entities"
	"liga-app/utils"

	"github.com/google/uuid"
	"google.golang.org/api/sheets/v4"
)

type API struct {
	srv *sheets.Service
}

func NewAPI(srv *sheets.Service) *API {
	return &API{srv: srv}
}

func (a *API) GetDrawings(leagueName string) ([]entities.GroupTeam, error) {
	spreadsheetId := "1OT-r1U4_xeRmg53Wj62ujbfto0aaURXynyu80ka2vOQ"

	if utils.IsDevEnv() {
		spreadsheetId = "16iVddhPYqP9oYhXh5Omt9YeAraoKPnkhK2suHtMSaHE"
	}

	readRange := fmt.Sprintf("%s!A1:E", leagueName)
	resp, err := a.srv.Spreadsheets.Values.Get(spreadsheetId, readRange).MajorDimension("ROWS").Do()
	if err != nil {
		return nil, fmt.Errorf("a.srv.Spreadsheets.Values.Get(spreadsheetId, readRange).MajorDimension(\"ROWS\").Do() : %w", err)
	}

	groupTeams := []entities.GroupTeam{}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for i, row := range resp.Values {
			if i == 0 {
				continue
			}
			// fmt.Printf("%s, %s, %s\n", row[0], row[1], row[2])
			groupName := row[0].(string)
			teamUUIDStr := row[2].(string)
			// groupName := row[1].(string)
			// teamName := row[2].(string)
			// leagueUUIDStr := row[4].(string)

			teamUUID, err := uuid.Parse(teamUUIDStr)
			if err != nil {
				return nil, fmt.Errorf("uuid.Parse(teamUUIDStr=%s): %w", teamUUIDStr, err)
			}

			gt := entities.GroupTeam{
				TeamUUID: teamUUID,
			}
			gt.Group.Name = groupName

			groupTeams = append(groupTeams, gt)
		}
	}

	return groupTeams, nil
}
