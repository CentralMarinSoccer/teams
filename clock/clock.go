package clock

import "time"

// Interface provides a mocking point for testing interactions with time
type Interface interface {
	Now() time.Time
}

// RealClock provides a real implementation using the time package
type RealClock struct{}

// Now returns the current time
func (RealClock) Now() time.Time { return time.Now() }

