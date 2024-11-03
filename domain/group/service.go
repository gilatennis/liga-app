package group

import (
	"context"
	"errors"
	"fmt"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
	"sort"

	"github.com/google/uuid"
)

type Service struct {
	groupRepo   repositories.IGroupRepository
	fixtureRepo repositories.IFixtureRepository
	matchRepo   repositories.IMatchRepository
}

func NewService(
	groupRepo repositories.IGroupRepository,
	fixtureRepo repositories.IFixtureRepository,
	matchRepo repositories.IMatchRepository,
) *Service {
	return &Service{
		groupRepo:   groupRepo,
		fixtureRepo: fixtureRepo,
		matchRepo:   matchRepo,
	}
}

func (s *Service) GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*entities.Group, error) {
	group, err := s.groupRepo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("s.groupRepo.GetByUUID(ctx, groupUUID=%s) : %w", groupUUID, err)
	}

	return group, nil
}

func (s *Service) Create(ctx context.Context, group *entities.Group) (*uuid.UUID, error) {
	groupUUID, err := s.groupRepo.Create(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("s.groupRepo.Create(ctx, group) : %w", err)
	}

	return groupUUID, nil
}

func (s *Service) AddTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error) {
	groupTeamUUID, err := s.groupRepo.AddTeam(ctx, groupUUID, teamUUID)
	if err != nil {
		return nil, fmt.Errorf("s.groupRepo.AddTeam(ctx, groupUUID=%s, teamUUID=%s) : %w", groupUUID, teamUUID, err)
	}

	return groupTeamUUID, nil
}

func (s *Service) RemoveTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error) {
	groupTeamUUID, err := s.groupRepo.RemoveTeam(ctx, groupUUID, teamUUID)
	if err != nil {
		return nil, fmt.Errorf("s.groupRepo.RemoveTeam(ctx, groupUUID=%s, teamUUID=%s) : %w", groupUUID, teamUUID, err)
	}

	return groupTeamUUID, nil
}

func (s *Service) GenerateGroupResults(ctx context.Context, groupUUID uuid.UUID) error {
	group, err := s.groupRepo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return fmt.Errorf("s.groupRepo.GetByUUID(ctx, groupUUID): %w", err)
	}

	teams, err := s.groupRepo.GetGroupTeams(ctx, groupUUID)
	if err != nil {
		return fmt.Errorf("s.groupRepo.GetGroupTeams(ctx, groupUUID): %w", err)
	}

	teamsByUUID := make(map[uuid.UUID]entities.Team)
	for _, groupTeam := range teams {
		teamsByUUID[groupTeam.UUID] = groupTeam
	}

	fixture, err := s.fixtureRepo.GetByUUID(ctx, group.FixtureUUID)
	if err != nil {
		return fmt.Errorf("s.fixtureRepo.GetByUUID(ctx, group.FixtureUUID) %w", err)
	}

	if fixture.Format == enums.FIXTURE_FORMAT_ROUND_ROBIN {
		teamResultByTeamUUID := make(map[uuid.UUID]*entities.GroupRoundRobinTeamResult)
		teamHeadToHeadResultByTeamUUID := make(map[uuid.UUID]map[uuid.UUID]*entities.GroupRoundRobinHeadToHeadResult)

		for _, t := range teams {
			teamResultByTeamUUID[t.UUID] = &entities.GroupRoundRobinTeamResult{
				Rank:        0,
				GroupUUID:   groupUUID,
				TeamUUID:    t.UUID,
				Players:     t.Players,
				Score:       0,
				MatchPlayed: 0,
				MatchWon:    0,
				MatchLost:   0,
			}
		}

		matches, _, err := s.matchRepo.Get(
			ctx,
			&common.Pagination{Page: 1, Limit: 100},
			&repositories.MatchFilter{GroupUUID: &groupUUID},
		)
		if err != nil {
			return fmt.Errorf("s.matchRepo.Get(...): %w", err)
		}

		for _, match := range matches {
			fmt.Println("match.Status", match.Status)
			if match.Status != string(enums.MATCH_STATUS_COMPLETED) {
				continue
			}

			teamResultByTeamUUID[*match.TeamAUUID].MatchPlayed += 1
			teamResultByTeamUUID[*match.TeamBUUID].MatchPlayed += 1

			// INFO: Playing a match adds 1 to score
			teamResultByTeamUUID[*match.TeamAUUID].Score += 1
			teamResultByTeamUUID[*match.TeamBUUID].Score += 1

			if teamHeadToHeadResultByTeamUUID[*match.TeamAUUID] == nil {
				teamHeadToHeadResultByTeamUUID[*match.TeamAUUID] = make(map[uuid.UUID]*entities.GroupRoundRobinHeadToHeadResult)
			}

			if teamHeadToHeadResultByTeamUUID[*match.TeamBUUID] == nil {
				teamHeadToHeadResultByTeamUUID[*match.TeamBUUID] = make(map[uuid.UUID]*entities.GroupRoundRobinHeadToHeadResult)
			}

			teamASetsWon := 0
			teamBSetsWon := 0
			for _, s := range match.Sets {
				if s.TeamAScore > s.TeamBScore {
					teamASetsWon += 1
				}
				if s.TeamBScore > s.TeamAScore {
					teamBSetsWon += 1
				}
			}

			if teamASetsWon > teamBSetsWon {
				// INFO: winning a match adds 1 to score
				teamResultByTeamUUID[*match.TeamAUUID].Score += 1

				teamResultByTeamUUID[*match.TeamAUUID].MatchWon += 1
				teamResultByTeamUUID[*match.TeamBUUID].MatchLost += 1
				// team A win
				teamHeadToHeadResultByTeamUUID[*match.TeamAUUID][*match.TeamBUUID] = &entities.GroupRoundRobinHeadToHeadResult{
					GroupUUID:        groupUUID,
					TeamUUID:         *match.TeamAUUID,
					TeamOpponentUUID: *match.TeamBUUID,
					Result:           enums.MatchResultWin,
				}
				// team B loss
				teamHeadToHeadResultByTeamUUID[*match.TeamBUUID][*match.TeamAUUID] = &entities.GroupRoundRobinHeadToHeadResult{
					GroupUUID:        groupUUID,
					TeamUUID:         *match.TeamBUUID,
					TeamOpponentUUID: *match.TeamBUUID,
					Result:           enums.MatchResultLoss,
				}
			}

			if teamBSetsWon > teamASetsWon {
				// INFO: winning a match adds 1 to score
				teamResultByTeamUUID[*match.TeamBUUID].Score += 1

				teamResultByTeamUUID[*match.TeamBUUID].MatchWon += 1
				teamResultByTeamUUID[*match.TeamAUUID].MatchLost += 1
				// team B win
				teamHeadToHeadResultByTeamUUID[*match.TeamBUUID][*match.TeamAUUID] = &entities.GroupRoundRobinHeadToHeadResult{
					GroupUUID:        groupUUID,
					TeamUUID:         *match.TeamBUUID,
					TeamOpponentUUID: *match.TeamAUUID,
					Result:           enums.MatchResultWin,
				}
				// team A loss
				teamHeadToHeadResultByTeamUUID[*match.TeamAUUID][*match.TeamBUUID] = &entities.GroupRoundRobinHeadToHeadResult{
					GroupUUID:        groupUUID,
					TeamUUID:         *match.TeamAUUID,
					TeamOpponentUUID: *match.TeamBUUID,
					Result:           enums.MatchResultLoss,
				}
			}
		}

		groupRoundRobinResults := make([]entities.GroupRoundRobinTeamResult, 0)
		groupRoundRobinHeadToHeadResults := make([]entities.GroupRoundRobinHeadToHeadResult, 0)

		for _, teamResult := range teamResultByTeamUUID {
			groupRoundRobinResults = append(groupRoundRobinResults, *teamResult)
		}

		for _, teamResultByOpponentUUID := range teamHeadToHeadResultByTeamUUID {
			for _, teamResult := range teamResultByOpponentUUID {
				groupRoundRobinHeadToHeadResults = append(groupRoundRobinHeadToHeadResults, *teamResult)
			}
		}

		sort.Slice(groupRoundRobinResults, func(i, j int) bool {
			if groupRoundRobinResults[i].Score == groupRoundRobinResults[j].Score {
				cumulativeGamesWon1 := groupRoundRobinResults[i].MatchWon - groupRoundRobinResults[i].MatchLost
				cumulativeGamesWon2 := groupRoundRobinResults[j].MatchWon - groupRoundRobinResults[j].MatchLost

				if cumulativeGamesWon1 == 0 && cumulativeGamesWon2 == 0 {
					return false
				}

				team1UUID := groupRoundRobinResults[i].TeamUUID
				team2UUID := groupRoundRobinResults[j].TeamUUID
				if cumulativeGamesWon1 == cumulativeGamesWon2 {
					if teamHeadToHeadResultByTeamUUID[team1UUID][team2UUID].Result == enums.MatchResultWin {
						return true
					} else {
						return false
					}
				} else {
					return cumulativeGamesWon1 > cumulativeGamesWon2
				}
			}

			return groupRoundRobinResults[i].Score > groupRoundRobinResults[j].Score
		})

		for i, res := range groupRoundRobinResults {
			groupRoundRobinResults[i].Rank = i + 1
			fmt.Printf("%s & %s - rank: %d matchPlayed: %d\n", res.Players[0].Name, res.Players[1].Name, i+1, res.MatchPlayed)
		}

		/*
			err = s.groupRepo.SaveGroupRoundRobinResults(ctx, groupUUID, groupRoundRobinResults, groupRoundRobinHeadToHeadResults)
			if err != nil {
				return fmt.Errorf("s.groupRepo.SaveGroupRoundRobinResults(ctx, groupUUID, groupRoundRobinResults, groupRoundRobinHeadToHeadResults): %w", err)
			}
		*/

		return nil
	}

	/**
		 * INFO: knockout result does not need to be generated
		 * It is easily calculated by looking at the last round of the fixture
	   * To get knockout result just use GetGroupResult method without generating beforehand
		 * */
	if fixture.Format == enums.FIXTURE_FORMAT_KNOCKOUT {
		return errors.New("STAGE_FORMAT_KNOCKOUT not implemented")
	}

	return nil
}
