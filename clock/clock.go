package clock

import "time"

type ClockInterface interface {
	Now() time.Time
	//	After(d time.Duration) <-chan time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

//func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
