package utils

import (
	"time"
)

func DisplayDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format(time.RFC3339)
}

func DisplayDateHuman(date time.Time) string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	if date.IsZero() {
		return ""
	}

	dateInLoc := date.In(loc)

	return dateInLoc.Format("2 January 2006")
}

func DisplayDateTimeHuman(dt time.Time) string {
	if dt.IsZero() {
		return ""
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return dt.In(loc).Format(time.RFC822)
}
