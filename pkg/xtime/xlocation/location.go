package xlocation

import "time"

// Must returns location according to specified time zone or causes a panic.
func Must(timeZone string) *time.Location {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		panic(err)
	}
	return loc
}
