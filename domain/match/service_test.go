package match_test

import (
	"fmt"
	"testing"
	"time"
)

func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	return date
}

func GetCurrentWeek(timezone *time.Location) int {
	now := time.Now().In(timezone)
	_, wkNum := now.ISOWeek()

	fmt.Println("wkNum", wkNum)

	return wkNum
}

func TestCreateMatchFlow(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	weekNum := GetCurrentWeek(loc)
	monday := FirstDayOfISOWeek(2024, weekNum, loc)
	saturday := monday.AddDate(0, 0, 5)
	sunday := monday.AddDate(0, 0, 6)

	fmt.Println("monday", monday)
	fmt.Println("saturday", saturday)
	fmt.Println("sunday", sunday)
}
