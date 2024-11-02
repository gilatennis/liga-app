package court

import (
	"context"
	"fmt"
	"liga-app/entities"
	"liga-app/interfaces/common"
	"liga-app/interfaces/repositories"
	"liga-app/utils"
	"sort"
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

func (s *Service) Get(c context.Context, p *common.Pagination) ([]entities.Court, error) {
	courts, err := s.courtRepo.Get(c, p)
	if err != nil {
		return nil, fmt.Errorf("s.courtRepo.Get(c, p) : %w", err)
	}

	return courts, nil
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

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, fmt.Errorf("time.LoadLocation(\"Asia/Jakarta\")")
	}

	now := time.Now().In(loc)
	beginningOfToday := utils.GetBeginningOfDay(now)
	beginningOfTomorrow := beginningOfToday.Add(24 * time.Hour)
	endOfLeague := time.Date(2025, 1, 15, 23, 59, 59, 59, now.Location())

	fmt.Println("beginningOfToday", beginningOfToday)
	fmt.Println("beginningOfTomorrow", beginningOfTomorrow)

	/*
		startOfWeek := s.getStartOfWeek()
		sixDays := 6 * 24 * time.Hour
		twentyThreeHours := (23 * time.Hour) + (59 * time.Minute) + (59 * time.Second)
		endOfWeek := startOfWeek.Add(sixDays).Add(twentyThreeHours)
	*/

	timeslots, err := s.courtTimeSlotRepo.GetByCourtSlug(c, court.Slug, beginningOfTomorrow, endOfLeague)
	if err != nil {
		return nil, fmt.Errorf("s.courtTimeSlotRepo.GetByCourtSlug(c, court.Slug, *rangeStart, *rangeEnd): %w", err)
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
	sort.Slice(daysSlice, func(i, j int) bool {
		return daysSlice[i].Before(daysSlice[j])
	})

	return daysSlice, nil
}

func (s *Service) GetAdminCourtAvailableDays(c context.Context, courtUUID uuid.UUID) ([]time.Time, error) {
	court, err := s.courtRepo.GetByUUID(c, courtUUID)
	if err != nil {
		return nil, fmt.Errorf("s.courtRepo.GetByUUID(c, courtUUID): %w", err)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, fmt.Errorf("time.LoadLocation(\"Asia/Jakarta\")")
	}

	now := time.Now().In(loc)
	endOfLeague := time.Date(2025, 1, 15, 23, 59, 59, 59, now.Location())
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, now.Location())

	timeslots, err := s.courtTimeSlotRepo.GetByCourtSlug(c, court.Slug, start, endOfLeague)
	if err != nil {
		return nil, fmt.Errorf("s.courtTimeSlotRepo.GetByCourtSlug(c, court.Slug, *rangeStart, *rangeEnd): %w", err)
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
	sort.Slice(daysSlice, func(i, j int) bool {
		return daysSlice[i].Before(daysSlice[j])
	})

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
