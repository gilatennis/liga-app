package group

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/repositories"

	"github.com/google/uuid"
)

type Service struct {
	groupRepo repositories.IGroupRepository
}

func NewService(groupRepo repositories.IGroupRepository) *Service {
	return &Service{groupRepo: groupRepo}
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
