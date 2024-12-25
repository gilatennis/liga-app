package team

import (
	"context"
	"fmt"
	"liga-app/config"
	domainErrors "liga-app/domain/errors"
	"liga-app/entities"
	"liga-app/enums"
	"liga-app/interfaces"
	iCommon "liga-app/interfaces/common"
	iRepo "liga-app/interfaces/repositories"
	"liga-app/interfaces/services"
	"mime/multipart"
	"sort"

	"github.com/google/uuid"
)

type Service struct {
	cfg         *config.Config
	teamRepo    iRepo.ITeamRepository
	groupRepo   iRepo.IGroupRepository
	fixtureRepo iRepo.IFixtureRepository
	fileStorage interfaces.IFileStorage
}

func NewService(
	cfg *config.Config,
	teamRepo iRepo.ITeamRepository,
	groupRepo iRepo.IGroupRepository,
	fixtureRepo iRepo.IFixtureRepository,
	fileStorage interfaces.IFileStorage,
) *Service {
	return &Service{
		cfg:         cfg,
		teamRepo:    teamRepo,
		groupRepo:   groupRepo,
		fixtureRepo: fixtureRepo,
		fileStorage: fileStorage,
	}
}

func (s *Service) Get(ctx context.Context, p *iCommon.Pagination, f *iRepo.TeamFilter, opt *services.TeamOptions) ([]entities.Team, int64, error) {
	teams, count, err := s.teamRepo.Get(ctx, p, f)
	if err != nil {
		return nil, 0, err
	}

	if opt != nil && opt.IncludeRoundRobinResult {
		format := enums.FIXTURE_FORMAT_ROUND_ROBIN
		ff := &iRepo.FixtureFilter{LeagueUUID: f.LeagueUUID, Format: &format}
		fixtures, count, err := s.fixtureRepo.Get(ctx, p, ff)
		if count > 1 || count == 0 {
			fmt.Printf("[ERROR] teamService: error in getting round robin result. fixtures.count should be 1 got %d instead \n", count)
			return teams, count, nil
		}
		if err != nil {
			return nil, 0, fmt.Errorf("s.fixtureRepo.Get(ctx, p, ff): %w", err)
		}
		roundRobinFixture := fixtures[0]

		for _, g := range roundRobinFixture.Groups {
			res, err := s.groupRepo.GetGroupRoundRobinResults(ctx, g.UUID)
			if err != nil {
				return nil, 0, fmt.Errorf("s.groupRepo.GetGroupRoundRobinResults(ctx, g.UUID=%s): %w", g.UUID, err)
			}

			for i, t := range teams {
				for _, r := range res {
					if r.TeamUUID == t.UUID {
						gr := g
						t.GroupStageGroup = &gr
						t.GroupStageRank = int64(r.Rank)
						teams[i] = t
					}
				}
			}
		}
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].GroupStageRank < teams[j].GroupStageRank
	})

	/*
		for _, t := range teams {
			fmt.Printf("Team: %s - Group: %s - Rank: %d \n", t.Name, t.GroupStageGroup.Name, t.GroupStageRank)
		}
	*/

	return teams, count, nil
}

func (s *Service) GetByUUID(ctx context.Context, teamUUID uuid.UUID) (*entities.Team, error) {
	team, err := s.teamRepo.GetByUUID(ctx, teamUUID)
	if err != nil {
		return nil, err
	}

	if team.ProfilePicture != "" {
		team.ProfilePictureURL, err = s.fileStorage.GetFileURL(ctx, s.cfg.Minio.Bucket, team.ProfilePicture)
		if err != nil {
			fmt.Printf("fileStorage.GetFileURL for team.ProfilePicture=%s failed: %v \n", team.ProfilePicture, err)
		}
	}

	return team, nil
}

func (s *Service) GetUserTeamInLeague(ctx context.Context, userUUID uuid.UUID, leagueUUID uuid.UUID) (*entities.Team, error) {
	p := iCommon.NewPagination(1, 1, []iCommon.OrderBy{{Column: "created_at", Direction: "desc"}})
	f := &iRepo.TeamFilter{LeagueUUID: &leagueUUID, PlayerUUID: &userUUID}
	teams, _, err := s.teamRepo.Get(ctx, p, f)
	if err != nil {
		return nil, err
	}

	if len(teams) == 0 {
		return nil, nil
	}

	team := teams[0]

	if team.ProfilePicture != "" {
		team.ProfilePictureURL, err = s.fileStorage.GetFileURL(ctx, s.cfg.Minio.Bucket, team.ProfilePicture)
		if err != nil {
			fmt.Printf("fileStorage.GetFileURL for team.ProfilePicture=%s failed: %v \n", team.ProfilePicture, err)
		}
	}

	return &team, nil
}

func (s *Service) GetLeagueTeamInviteLink(ctx context.Context, leagueUUID uuid.UUID, teamUUID uuid.UUID) string {
	inviteLink := fmt.Sprintf("%s/invite?league_uuid=%s&team_uuid=%s", s.cfg.App.BaseURL, leagueUUID, teamUUID)
	return inviteLink
}

func (s *Service) Update(ctx context.Context, user *entities.User, team *entities.Team, profilePicture *multipart.File) (*uuid.UUID, error) {
	if user.UUID != team.CaptainUserUUID {
		return nil, &domainErrors.ErrUnauthorized{Message: "unauthorized. only captain can update team info"}
	}

	if profilePicture != nil {
		profilePicName, err := s.fileStorage.CreateFileName(ctx, *profilePicture)
		if err != nil {
			return nil, err
		}

		err = s.fileStorage.SaveFile(ctx, s.cfg.Minio.Bucket, profilePicName, *profilePicture)
		if err != nil {
			return nil, err
		}

		team.ProfilePicture = profilePicName
	}

	userUUID, err := s.teamRepo.Update(ctx, team)
	if err != nil {
		return nil, err
	}

	return userUUID, nil
}
