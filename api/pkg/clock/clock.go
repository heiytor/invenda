package clock

import "time"

type Clock interface {
	Now() time.Time
}

// Now returns the current [time.Time] with the locaion set to UTC.
func Now() time.Time {
	return Backend.Now()
}
