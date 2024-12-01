package group

import (
	"context"
	"errors"
	"fmt"
	"liga-app/entities"
	"strings"

	domainErrors "liga-app/domain/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db}
}

func (s *Store) GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*entities.Group, error) {
	g := entities.Group{}
	err := s.db.
		Preload("Fixture").
		Preload("Fixture.League").
		Preload("Teams").
		Preload("Teams.Captain").
		Preload("Teams.Players").
		Preload("Teams.Players.Ratings").
		First(&g, "uuid = ?", groupUUID).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db.First(&g, \"uuid = ?\", groupUUID).Error : %w", err)
	}

	if len(g.Teams) == 0 {
		return &g, nil
	}

	teamUUIDs := []uuid.UUID{}
	for _, t := range g.Teams {
		teamUUIDs = append(teamUUIDs, t.UUID)
	}

	from := fmt.Sprintf(`FROM (
    SELECT '%s' AS team_uuid
  `, teamUUIDs[0].String())
	for _, u := range teamUUIDs[1:] {
		from += fmt.Sprintf("  UNION ALL SELECT '%s' \n", u.String())
	}
	from += ") AS tl"

	teamUUIDString := []string{}
	for _, u := range teamUUIDs {
		teamUUIDString = append(teamUUIDString, fmt.Sprintf("'%s'", u.String()))
	}
	teamUUIDStringJoin := strings.Join(teamUUIDString, ", ")
	order := fmt.Sprintf("ORDER BY FIELD(tl.team_uuid, %s)", teamUUIDStringJoin)

	rawQuery := fmt.Sprintf(`
    SELECT
        tl.team_uuid,
        (SELECT COUNT(DISTINCT m.uuid)
        FROM matches m
        JOIN rounds r ON m.round_uuid = r.uuid
        JOIN `+"`groups`"+` g ON r.group_uuid = g.uuid
        JOIN fixtures f ON g.fixture_uuid = f.uuid
        WHERE (m.team_a_uuid = tl.team_uuid OR m.team_b_uuid = tl.team_uuid)
          AND f.format = 'ROUND_ROBIN' AND m.team_a_uuid IS NOT NULL AND m.team_b_uuid IS NOT NULL) AS group_stage_match_count,
        (SELECT COUNT(DISTINCT m.uuid)
        FROM matches m
        JOIN rounds r ON m.round_uuid = r.uuid
        JOIN `+"`groups`"+` g ON r.group_uuid = g.uuid
        JOIN fixtures f ON g.fixture_uuid = f.uuid
        WHERE (m.team_a_uuid = tl.team_uuid OR m.team_b_uuid = tl.team_uuid)
          AND f.format = 'ROUND_ROBIN' AND m.status = 'COMPLETED') AS group_stage_match_completed_count
      %s
      %s;
  `, from, order)

	type Temp struct {
		TeamUUID                      uuid.UUID
		GroupStageMatchCount          int64
		GroupStageMatchCompletedCount int64
	}
	temps := []Temp{}
	err = s.db.WithContext(ctx).
		Raw(rawQuery).
		Scan(&temps).
		Error
	if err != nil {
		return nil, err
	}

	for i, t := range g.Teams {
		t.GroupStageMatchCount = temps[i].GroupStageMatchCount
		t.GroupStageMatchCompletedCount = temps[i].GroupStageMatchCompletedCount
		g.Teams[i] = t
	}

	return &g, nil
}

func (s *Store) Create(ctx context.Context, group *entities.Group) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).Create(group).Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(ctx).Create(group).Error : %w", err)
	}

	return &group.UUID, nil
}

func (s *Store) AddTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error) {
	existing := entities.GroupTeam{}
	err := s.db.WithContext(ctx).First(&existing, "group_uuid = ? AND team_uuid = ?", groupUUID, teamUUID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("s.db.First(...).Error :%w", err)
	}

	if existing.UUID != uuid.Nil {
		return nil, &domainErrors.ErrConflict{Message: "already in team"}
	}

	gt := entities.GroupTeam{
		GroupUUID: groupUUID,
		TeamUUID:  teamUUID,
	}

	err = s.db.WithContext(ctx).Create(&gt).Error
	if err != nil {
		return nil, fmt.Errorf("s.db.WithContext(ctx).Create(&gt).Error : %w", err)
	}

	return &gt.UUID, nil
}

func (s *Store) RemoveTeam(ctx context.Context, groupUUID uuid.UUID, teamUUID uuid.UUID) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Unscoped().
		Delete(&entities.GroupTeam{}, "group_uuid = ? AND team_uuid = ?", groupUUID, teamUUID).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db...Delete(&entities.GroupTeam{}, \"group_uuid = ? AND team_uuid = ?\", groupUUID=%s, teamUUID=%s).Error : %w", groupUUID, teamUUID, err)
	}
	return nil, nil
}

func (s *Store) GetGroupTeams(ctx context.Context, groupUUID uuid.UUID) ([]entities.Team, error) {
	teams := []entities.Team{}
	err := s.db.WithContext(ctx).
		Preload("Captain").
		Preload("Players").
		Joins("JOIN group_teams on group_teams.team_uuid = teams.uuid").
		Find(&teams, "group_teams.group_uuid = ?", groupUUID).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db...Find(&teams, \"group_teams.group_uuid = ?\", groupUUID=%s) : %w", groupUUID, err)
	}

	return teams, nil
}

func (s *Store) SaveGroupRoundRobinResults(
	ctx context.Context,
	groupUUID uuid.UUID,
	groupRoundRobinResults []entities.GroupRoundRobinTeamResult,
	groupRoundRobinHeadToHeadResults []entities.GroupRoundRobinHeadToHeadResult,
) error {
	tx := s.db.WithContext(ctx).Begin()

	err := tx.Unscoped().Delete(&entities.GroupRoundRobinTeamResult{}, "group_uuid = ?", groupUUID).Error
	if err != nil {
		return fmt.Errorf("tx.Unscoped().Delete(&entities.GroupRoundRobinTeamResult{}, \"group_uuid = ?\", groupUUID).Error : %w", err)
	}

	err = tx.Unscoped().Delete(&entities.GroupRoundRobinHeadToHeadResult{}, "group_uuid = ?", groupUUID).Error
	if err != nil {
		return fmt.Errorf("tx.Unscoped().Delete(&entities.GroupRoundRobinHeadToHeadResult{}, \"group_uuid = ?\", groupUUID).Error : %w", err)
	}

	err = tx.Save(groupRoundRobinResults).Error
	if err != nil {
		return fmt.Errorf("tx.Save(groupRoundRobinResults).Error : %w", err)
	}

	if len(groupRoundRobinHeadToHeadResults) > 0 {
		err = tx.Save(groupRoundRobinHeadToHeadResults).Error
		if err != nil {
			return fmt.Errorf("tx.Save(groupRoundRobinHeadToHeadResults).Error : %w", err)
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("tx.Commit().Error : %w", err)
	}

	return nil
}

func (s *Store) GetGroupRoundRobinResults(ctx context.Context, groupUUID uuid.UUID) ([]entities.GroupRoundRobinTeamResult, error) {
	results := []entities.GroupRoundRobinTeamResult{}

	err := s.db.WithContext(ctx).
		Debug().
		Preload("Group").
		Preload("Team").
		Where("group_uuid = ?", groupUUID).
		Order("`rank` ASC").
		Find(&results).
		Error
	if err != nil {
		return nil, fmt.Errorf("s.db...Find(&results, \"group_uuid = %s): %w\"", groupUUID, err)
	}

	return results, nil
}
