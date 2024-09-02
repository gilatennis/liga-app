package league_registration

import (
	"context"
	"liga-app/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db}
}

func (s *Store) GetByUUID(ctx context.Context, leagueRegUUID uuid.UUID) (*entities.LeagueRegistration, error) {
	var leagueReg entities.LeagueRegistration
	err := s.db.WithContext(ctx).
		Where("uuid = ?", leagueRegUUID).
		Preload("League").
		Preload("User").
		First(&leagueReg).
		Error
	if err != nil {
		return nil, err
	}

	return &leagueReg, nil
}

func (s *Store) Update(ctx context.Context, leagueReg *entities.LeagueRegistration) (*uuid.UUID, error) {
	err := s.db.WithContext(ctx).
		Omit(clause.Associations).
		Updates(leagueReg).
		Error
	if err != nil {
		return nil, err
	}

	return &leagueReg.UUID, nil
}

func (s *Store) CompleteRegistration(ctx context.Context, leagueReg *entities.LeagueRegistration) (*uuid.UUID, error) {
	tx := s.db.WithContext(ctx).Begin()

	err := tx.
		Omit(clause.Associations).
		Updates(leagueReg).
		Error
	if err != nil {
		return nil, err
	}

	team := entities.Team{
		Name:            leagueReg.User.Name,
		CaptainUserUUID: leagueReg.UserUUID,
	}
	err = tx.Omit(clause.Associations).Save(&team).Error
	if err != nil {
		return nil, err
	}

	leagueTeam := entities.LeagueTeam{
		LeagueUUID: leagueReg.LeagueUUID,
		TeamUUID:   team.UUID,
	}
	err = tx.Omit(clause.Associations).Save(&leagueTeam).Error
	if err != nil {
		return nil, err
	}

	teamPlayer := entities.TeamPlayer{
		TeamUUID:   team.UUID,
		PlayerUUID: leagueReg.UserUUID,
	}
	err = tx.Omit(clause.Associations).Save(&teamPlayer).Error
	if err != nil {
		return nil, err
	}

	leaguePlayer := entities.LeaguePlayer{
		LeagueUUID: leagueReg.LeagueUUID,
		PlayerUUID: leagueReg.UserUUID,
	}
	err = tx.Omit(clause.Associations).Save(&leaguePlayer).Error
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &leagueReg.UUID, nil
}
