package clock

import "time"

type backend struct{}

var Backend Clock = (*backend)(nil)

func (*backend) Now() time.Time {
	return time.Now().UTC()
}
