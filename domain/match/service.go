package match

import (
	iRepo "liga-app/interfaces/repositories"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	courtTimeSlotStore iRepo.ICourtTimeSlotRepository
}

func NewService(courtTimeSlotStore iRepo.ICourtTimeSlotRepository) *Service {
	return &Service{
		courtTimeSlotStore: courtTimeSlotStore,
	}
}

func (s *Service) ProposeMatchSchedule(userUUID uuid.UUID, start time.Time, end time.Time) error {
	return nil
}

func (s *Service) AcceptMatchSchedule(matchScheduleRequestUUID uuid.UUID) error {
	return nil
}
