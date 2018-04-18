package utils

import "time"

func ParseTime(timeRaw string) time.Time {
	const timeLayout = "02-01-2006 15:04 (MST)"
	t, _ := time.Parse(timeLayout, timeRaw)
	return t
}
