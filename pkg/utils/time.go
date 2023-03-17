package utils

import "time"

func UpdateTimeZone(inTime time.Time, loc *time.Location) time.Time {
	return time.Date(
		inTime.Year(),
		inTime.Month(),
		inTime.Day(),
		inTime.Hour(),
		inTime.Minute(),
		inTime.Second(),
		inTime.Nanosecond(),
		loc,
	)
}
