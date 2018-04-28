package utils

import "time"

var result Result

func ParseTime(timeRaw string) time.Time {
	const timeLayout = "02-01-2006 15:04 (MST)"
	t, _ := time.Parse(timeLayout, timeRaw)
	return t
}
