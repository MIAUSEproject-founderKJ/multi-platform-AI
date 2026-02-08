//internal/scheduler/ticker.go

package scheduler

import "time"

// NewPreciseTicker shims the standard ticker for now.
// In high-performance scenarios, this would use syscall.Timer or similar.
func NewPreciseTicker(hz int) *time.Ticker {
	if hz <= 0 {
		hz = 10 // Default to 10Hz safety floor
	}
	interval := time.Second / time.Duration(hz)
	return time.NewTicker(interval)
}
