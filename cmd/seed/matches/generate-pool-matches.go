package matches

import (
	"context"
	"errors"
	"fmt"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/infra/gsheet"
	mysqlLeague "liga-app/infra/mysql/league"
	mysqlTeam "liga-app/infra/mysql/team"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GeneratePoolMatches(db *gorm.DB, gsheetAPI *gsheet.API) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		leagueStore := mysqlLeague.NewStore(db)
		teamStore := mysqlTeam.NewStore(db)
		_ = teamStore

		// 3c508a53-b399-46b9-b5c5-75b12e5092fc
		leagueUUIDStr := c.String("league-uuid")
		leagueUUID, err := uuid.Parse(leagueUUIDStr)
		if err != nil {
			return fmt.Errorf("uuid.Parse(leagueUUIDStr=%s): %w", leagueUUIDStr, err)
		}

		league, err := leagueStore.GetByUUID(c.Context, leagueUUID)
		if err != nil {
			return fmt.Errorf("leagueStore.GetByUUID(...): %w", err)
		}

		// delete matches -> rounds -> groups -> fixtures for the league
		err = DeleteLeagueRoundRobinFixture(c.Context, db, leagueUUID)
		if err != nil {
			return fmt.Errorf("DeleteLeagueRoundRobinFixture(c.Context, db, leagueUUID=%s): %w", leagueUUID, err)
		}

		fixture := entities.Fixture{
			LeagueUUID:  leagueUUID,
			Name:        "Round Robin",
			Format:      enums.FIXTURE_FORMAT_ROUND_ROBIN,
			StageNumber: 1,
		}
		err = db.Create(&fixture).Error
		if err != nil {
			return fmt.Errorf("db.Create(&fixture).Error")
		}

		groupTeamsFromGSheet, err := gsheetAPI.Drawings.GetDrawings(league.Name)
		if err != nil {
			return fmt.Errorf("gsheetAPI.Drawings.GetDrawings(league.Name=%s): %w", league.Name, err)
		}

		groups, err := CreateGroupsAndAssignTeams(db, fixture.UUID, groupTeamsFromGSheet)
		if err != nil {
			return fmt.Errorf("CreateGroupsAndAssignTeams(db, fixture.UUID, groupTeams): %w", err)
		}

		for _, g := range groups {
			// pp.Println("g.Teams", g.Teams)

			teamCount := len(g.Teams)

			teamIDs := make([]string, teamCount)
			for i, team := range g.Teams {
				teamIDs[i] = team.UUID.String()
			}

			if len(teamIDs) < teamCount {
				for i := len(teamIDs); i < teamCount; i++ {
					teamIDs = append(teamIDs, "TBA")
				}
			}

			rounds, err := GenerateRoundRobinMatches(g.UUID, teamIDs)
			if err != nil {
				return fmt.Errorf("GenerateRoundRobinMatches(groupUUID, teamIDs) : %w", err)
			}
			_ = rounds

			// pp.Println("rounds", rounds)
			PrintRounds(rounds)

			err = SaveRoundsToDB(c.Context, db, rounds, g.UUID, leagueUUID)
			if err != nil {
				return fmt.Errorf("SaveRoundsToDB(c.Context, db, rounds, g.UUID, leagueUUID) : %w", err)
			}
		}

		return nil
	}
}

func PrintRounds(rounds []entities.Round) {
	for _, r := range rounds {
		fmt.Printf("Round number: %d \n", r.Number)
		for _, m := range r.Matches {
			fmt.Printf("\t Match Number: %d \n", m.Number)
			fmt.Printf("\t\t Team A: %s", m.TeamAUUID)
			fmt.Printf("\t\t Team B: %s", m.TeamBUUID)
			fmt.Println()
		}
		fmt.Println()
	}
}

func GenerateRoundRobinMatches(groupUUID uuid.UUID, teamIDs []string) ([]entities.Round, error) {
	// assign number to teamIDs
	numToTeamID := make(map[int]string)
	var teamNums []int
	const BYE = "BYE"

	if len(teamIDs)%2 != 0 {
		teamNums = make([]int, len(teamIDs)+1)
		teamIDs = append(teamIDs, BYE)
	} else {
		teamNums = make([]int, len(teamIDs))
	}

	for i, teamID := range teamIDs {
		numToTeamID[i] = teamID
		teamNums[i] = i
	}

	numRounds := len(teamIDs) - 1
	rounds := []entities.Round{}
	for i := 0; i < numRounds; i++ {
		r := entities.Round{
			Base: entities.Base{
				UUID: uuid.New(),
			},
			GroupUUID: groupUUID,
			Number:    i + 1,
		}

		for j := 0; j < len(teamNums)/2; j++ {
			left := teamNums[j]
			right := teamNums[len(teamNums)-j-1]

			opp1ID := numToTeamID[left]
			opp2ID := numToTeamID[right]

			match := entities.Match{
				Base: entities.Base{
					UUID: uuid.New(),
				},
				Number:    j + 1,
				RoundUUID: r.UUID,
				// RoundNumber: r.Number,
				Status: string(enums.MATCH_STATUS_WAITING),
			}

			oppIDs := []string{opp1ID, opp2ID}
			for i := 0; i < 2; i++ {
				oppID := oppIDs[i]

				if i == 0 {
					if oppID == BYE {
						match.IsByeTeamA = true
						continue
					}

					if oppID == "TBA" {
						continue
					}

					oppUUID, err := uuid.Parse(oppID)
					if err != nil {
						err = fmt.Errorf("uuid.Parse(oppID=%s): %w", oppID, err)
						return nil, err
					}

					match.TeamAUUID = &oppUUID
				}

				if i == 1 {
					if oppID == BYE {
						match.IsByeTeamB = true
						continue
					}

					if oppID == "TBA" {
						continue
					}

					oppUUID, err := uuid.Parse(oppID)
					if err != nil {
						err = fmt.Errorf("uuid.Parse(oppID=%s): %w", oppID, err)
						return nil, err
					}

					match.TeamBUUID = &oppUUID
				}
			}

			r.Matches = append(r.Matches, match)
		}

		rounds = append(rounds, r)

		if i == numRounds-1 {
			break
		}
		// rotate teamIDs
		rTemp := 0
		for k := 0; k < len(teamNums); k++ {
			// keep first idx
			if k == 0 {
				continue
			}
			if k == 1 {
				rTemp = teamNums[k]
				teamNums[k] = teamNums[len(teamNums)-1]
				continue
			}
			temp := teamNums[k]
			teamNums[k] = rTemp
			rTemp = temp
		}
	}

	return rounds, nil
}

func DeleteLeagueRoundRobinFixture(ctx context.Context, db *gorm.DB, leagueUUID uuid.UUID) error {
	tx := db.WithContext(ctx).Session(&gorm.Session{}).Begin()

	league := entities.League{}
	err := tx.Where("uuid = ?", leagueUUID).First(&league).Error
	if err != nil {
		return fmt.Errorf("db.Where(\"uuid = ?\", leagueUUID=%s).First(&league).Error : %w", leagueUUID, err)
	}

	fixtureRoundRobin := entities.Fixture{}
	err = db.
		Where("league_uuid = ?", league.UUID).
		Where("format = ?", enums.FIXTURE_FORMAT_ROUND_ROBIN).
		First(&fixtureRoundRobin).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("db.Where(...).Where(...).First(&fixtureRoundRobin) : %w", err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	groups := []entities.Group{}
	err = tx.Where("fixture_uuid = ?", fixtureRoundRobin.UUID).Find(&groups).Error
	if err != nil {
		return fmt.Errorf("tx.Where(\"fixture_uuid = ?\", f.UUID).Find(&groups).Error : %w", err)
	}

	for _, g := range groups {
		err = tx.Unscoped().Delete(&entities.GroupTeam{}, "group_uuid = ?", g.UUID).Error
		if err != nil {
			return fmt.Errorf("tx.Unscoped().Delete(&entities.GroupTeam{}, \"group_uuid = ?\", g.UUID).Error : %w", err)
		}

		rounds := []entities.Round{}
		err = tx.Where("group_uuid = ?", g.UUID).Find(&rounds).Error
		if err != nil {
			return fmt.Errorf("tx.Where(\"group_uuid = ?\", f.UUID).Find(&rounds).Error: %w", err)
		}

		for _, r := range rounds {
			matches := []entities.Match{}
			err = tx.Where("round_uuid = ?", r.UUID).Find(&matches).Error
			if err != nil {
				return fmt.Errorf("tx.Where(\"round_uuid = ?\", r.UUID).Find(&matches).Error : %w", err)
			}

			for _, m := range matches {
				// INFO: Delete Match Sets
				err = tx.Unscoped().Delete(&entities.MatchSet{}, "match_uuid = ?", m.UUID).Error
				if err != nil {
					return fmt.Errorf("tx.Unscoped().Delete(&entities.MatchSet{}, \"match_uuid = ?\", m.UUID).Error : %w", err)
				}
			}

			// INFO: Delete Matches
			err = tx.Unscoped().Delete(&entities.Match{}, "round_uuid = ?", r.UUID).Error
			if err != nil {
				return fmt.Errorf("tx.Unscoped().Delete(&entities.MatchSet{}, \"match_uuid = ?\", m.UUID).Error : %w", err)
			}
		}

		// INFO: Delete Rounds
		err = tx.Unscoped().Delete(&entities.Round{}, "group_uuid = ?", g.UUID).Error
		if err != nil {
			return fmt.Errorf("tx.Unscoped().Delete(&entities.Round{}, \"group_uuid = ?\", m.UUID).Error : %w", err)
		}
	}

	// INFO: Delete Groups
	err = tx.Unscoped().Delete(&entities.Group{}, "fixture_uuid = ?", fixtureRoundRobin.UUID).Error
	if err != nil {
		return fmt.Errorf("tx.Delete(&entities.Group{}, \"fixture_uuid = ?\", f.UUID).Error: %w", err)
	}

	// INFO: Delete Fixture
	err = tx.Unscoped().Delete(&entities.Fixture{}, "league_uuid = ? AND format = ?", leagueUUID, enums.FIXTURE_FORMAT_ROUND_ROBIN).Error
	if err != nil {
		return fmt.Errorf("db.Delete(&entities.Fixture{}).Where(\"league_uuid = ?\", leagueUUIDStr).Error")
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("tx.Commit().Error")
	}

	return nil
}

func CreateGroupsAndAssignTeams(db *gorm.DB, fixtureUUID uuid.UUID, groupTeams []entities.GroupTeam) ([]entities.Group, error) {
	mapGroupNameTeamUUIDs := make(map[string][]uuid.UUID)

	for _, gt := range groupTeams {
		if _, ok := mapGroupNameTeamUUIDs[gt.Group.Name]; !ok {
			mapGroupNameTeamUUIDs[gt.Group.Name] = []uuid.UUID{}
		}
		fmt.Printf("group_name = %s, team_uuid = %s \n", gt.Group.Name, gt.TeamUUID)
		teamUUIDs := mapGroupNameTeamUUIDs[gt.Group.Name]
		teamUUIDs = append(teamUUIDs, gt.TeamUUID)
		mapGroupNameTeamUUIDs[gt.Group.Name] = teamUUIDs
	}

	groups := []entities.Group{}
	groupTeams = []entities.GroupTeam{}
	for gName, teamUUIDs := range mapGroupNameTeamUUIDs {
		group := entities.Group{
			FixtureUUID: fixtureUUID,
			Name:        gName,
		}

		err := db.Create(&group).Error
		if err != nil {
			return nil, fmt.Errorf("db.Create(&group).Error")
		}

		groups = append(groups, group)

		for _, tUUID := range teamUUIDs {
			groupTeams = append(groupTeams, entities.GroupTeam{
				GroupUUID: group.UUID,
				TeamUUID:  tUUID,
			})
		}
	}

	err := db.Create(&groupTeams).Error
	if err != nil {
		return nil, fmt.Errorf("db.Create(&groupTeams).Error")
	}

	groups = []entities.Group{}
	err = db.Preload("Teams").Where("fixture_uuid = ?", fixtureUUID).Find(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("db.Where(\"fixture_uuid = ?\", fixtureUUID).Find(&groups).Error")
	}

	return groups, nil
}

func SaveRoundsToDB(
	ctx context.Context,
	db *gorm.DB,
	rounds []entities.Round,
	groupUUID uuid.UUID,
	leagueUUID uuid.UUID,
) error {
	tx := db.WithContext(ctx).Begin()

	for _, r := range rounds {
		r.GroupUUID = groupUUID

		err := tx.Omit(clause.Associations).Create(&r).Error
		if err != nil {
			return fmt.Errorf("tx.Omit(clause.Associations).Create(r).Error: %w", err)
		}

		for j, m := range r.Matches {
			m.LeagueUUID = leagueUUID
			m.RoundUUID = r.UUID

			r.Matches[j] = m
		}

		err = tx.Omit(clause.Associations).Create(r.Matches).Error
		if err != nil {
			return fmt.Errorf("tx.Omit(clause.Associations).Create(r.Matches).Error : %w", err)
		}
	}

	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("tx.Commit(): %w", err)
	}

	return nil
}
