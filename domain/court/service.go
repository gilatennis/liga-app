package court

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/repositories"
	"liga-app/utils"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	courtRepo         repositories.ICourtRepository
	courtTimeSlotRepo repositories.ICourtTimeSlotRepository
}

func NewService(
	courtRepo repositories.ICourtRepository,
	courtTimeSlotRepo repositories.ICourtTimeSlotRepository,
) *Service {
	return &Service{
		courtRepo:         courtRepo,
		courtTimeSlotRepo: courtTimeSlotRepo,
	}
}

func (s *Service) GetByUUID(c context.Context, courtUUID uuid.UUID) (*entities.Court, error) {
	court, err := s.courtRepo.GetByUUID(c, courtUUID)
	if err != nil {
		return nil, fmt.Errorf("s.courtRepo.GetByUUID(c, courtUUID)")
	}

	return court, nil
}

func (s *Service) getStartOfWeek() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	thisYear := utils.GetCurrentYear(loc)
	thisWeek := utils.GetCurrentWeek(loc)
	monday := utils.FirstDayOfISOWeek(thisYear, thisWeek, loc)
	return monday
}

func (s *Service) GetCourtAvailableDays(c context.Context, courtUUID uuid.UUID) ([]time.Time, error) {
	court, err := s.courtRepo.GetByUUID(c, courtUUID)
	if err != nil {
		return nil, fmt.Errorf("s.courtRepo.GetByUUID(c, courtUUID): %w", err)
	}

	startOfWeek := s.getStartOfWeek()
	sixDays := 6 * 24 * time.Hour
	twentyThreeHours := (23 * time.Hour) + (59 * time.Minute) + (59 * time.Second)
	endOfWeek := startOfWeek.Add(sixDays).Add(twentyThreeHours)

	timeslots, err := s.courtTimeSlotRepo.GetByCourtSlug(c, court.Slug, startOfWeek, endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("s.courtTimeSlotRepo.GetByCourtSlug(c, court.Slug, *rangeStart, *rangeEnd): %w", err)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, fmt.Errorf("time.LoadLocation(\"Asia/Jakarta\")")
	}

	daysMap := make(map[int]time.Time)
	for _, ts := range timeslots {
		day := ts.Start.Day()
		if _, ok := daysMap[day]; !ok {
			daysMap[day] = ts.Start.Truncate(24 * time.Hour).In(loc)
		}
	}

	daysSlice := []time.Time{}
	for _, d := range daysMap {
		daysSlice = append(daysSlice, d)
	}

	return daysSlice, nil
}

func (s *Service) GetCourtTimeSlotsByDay(c context.Context, courtUUID uuid.UUID, day time.Time) ([]entities.CourtTimeSlot, error) {
	court, err := s.courtRepo.GetByUUID(c, courtUUID)
	if err != nil {
		return nil, fmt.Errorf("s.courtRepo.GetByUUID(c, courtUUID): %w", err)
	}

	twentyThreeHours := (23 * time.Hour) + (59 * time.Minute) + (59 * time.Second)
	endOfDay := day.Add(twentyThreeHours)

	timeslots, err := s.courtTimeSlotRepo.GetByCourtSlug(c, court.Slug, day, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("s.courtTimeSlotRepo.GetByCourtSlug(c, court.Slug, day, endOfDay)")
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, fmt.Errorf("time.LoadLocation(\"Asia/Jakarta\")")
	}

	for i, ts := range timeslots {
		ts.Start = ts.Start.In(loc)
		ts.End = ts.End.In(loc)
		timeslots[i] = ts
	}

	return timeslots, nil
}
