package mocks

import (
	"time"
)

// ClockData provides a mechanism to set and get mocked ClockData
type ClockData struct {
	NowCall struct {
		Returns struct {
			Time time.Time
		}
	}
}

// Now provides a mechanism to receive consistent time / date data
func (ClockData) Now() time.Time {
	t1, _ := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	return t1
}
