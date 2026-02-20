package utime

import (
	"context"
	"sync/atomic"
	"time"
)

// Time represents a cached timestamp as nanoseconds since Unix epoch.
// It provides high-performance time operations with minimal overhead.
type Time int64

var Utime = NewClock(500 * time.Microsecond)

// Clock is an instance of a clock whose time increments roughly at a
// configured granularity, but lookups are effectively free relative to
// normal [time.Now]. The struct is cache line aligned for optimal performance.
type Clock struct {
	// Most frequently accessed field first for optimal cache utilization
	now    atomic.Int64       // 8 bytes
	ctx    context.Context    // 8 bytes (interface)
	cancel context.CancelFunc // 8 bytes (interface)
	_      [40]byte           // 64 - 8 - 8 - 8 = 40 bytes padding
}

// NewClock creates a new Clock configured to tick at approximately
// granularity intervals. Clock is running when created, and may be stopped
// by calling Clock.Stop. A stopped Clock cannot be resumed.
//
//go:inline
func NewClock(granularity time.Duration) *Clock {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Clock{ctx: ctx, cancel: cancel}
	c.now.Store(time.Now().UnixNano())
	go c.run(granularity)
	return c
}

// Now returns a Time that represents the current cached time.
// The Time returned will never be in the future, but will always be
// less than or equal to the actual current time.
//
//go:inline
func (c *Clock) Now() Time {
	return Time(c.now.Load())
}

// Unix returns the cached time as Unix timestamp in seconds.
//
//go:inline
func (t Time) Unix() int64 {
	return int64(t) / int64(time.Second)
}

// UnixMilli returns the cached time as Unix timestamp in milliseconds.
//
//go:inline
func (t Time) UnixMilli() int64 {
	return int64(t) / int64(time.Millisecond)
}

// UnixNano returns the cached time as Unix timestamp in nanoseconds.
//
//go:inline
func (t Time) UnixNano() int64 {
	return int64(t)
}

// ToTime converts the cached time to a time.Time value.
//
//go:inline
func (t Time) ToTime() time.Time {
	return time.Unix(0, int64(t))
}

// Add adds a duration to the cached time and returns the new Time.
//
//go:inline
func (t Time) Add(d time.Duration) Time {
	return Time(int64(t) + int64(d))
}

func (t Time) AddDate(years, months, days int) Time {
	// Convert to time.Time, add date, and convert back to Time
	tTime := t.ToTime().AddDate(years, months, days)
	return Time(tTime.UnixNano())
}

// Sub subtracts a duration from the cached time and returns the new Time.
//
//go:inline
func (t Time) Sub(d Time) time.Duration {
	return time.Duration(int64(t) - int64(d))
}

func Until(d Time) time.Duration {
	return time.Duration(int64(d) - int64(Utime.Now()))
}

// Since returns time.Duration since the given time relative to the Clock's
// current cached time.
//
//go:inline
func (t Time) Since(i Time) time.Duration {
	return time.Duration(int64(t) - int64(i))
}

func (t Time) Before(u Time) bool {
	return t.ToTime().Before(u.ToTime())
}

func (t Time) After(u Time) bool {
	return t.ToTime().After(u.ToTime())
}

func FromUnixNano(u int64) Time {
	return Time(u)
}

// Stop stops the Clock ticker and cannot be resumed.
func (c *Clock) Stop() {
	c.cancel()
}

// run is the internal goroutine that updates the cached time at the specified granularity.
//
//go:inline
func (c *Clock) run(granularity time.Duration) {
	t := time.NewTicker(granularity)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			c.now.Store(time.Now().UnixNano())
		case <-c.ctx.Done():
			return
		}
	}
}

// FromTime converts a time.Time value to a cached Time.
//
//go:inline
func (c *Clock) FromTime(t time.Time) Time {
	return Time(t.UnixNano())
}
