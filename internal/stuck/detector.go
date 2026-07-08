package stuck

import (
	"sync"
	"time"
)

// State represents the detector's current state.
type State int

const (
	StateIdle State = iota
	StateResponding
	StateStuck
)

// Detector watches for Claude responses that stop producing tokens.
type Detector struct {
	timeout time.Duration
	timer   *time.Timer
	state   State
	mu      sync.Mutex
	onStuck func()
}

// NewDetector creates a detector that calls onStuck after timeout passes
// without a token being received while responding.
func NewDetector(timeout time.Duration, onStuck func()) *Detector {
	return &Detector{
		timeout: timeout,
		state:   StateIdle,
		onStuck: onStuck,
	}
}

// OnResponseStart transitions the detector to the responding state.
func (d *Detector) OnResponseStart() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.state = StateResponding
	d.resetTimerLocked()
}

// OnTokenReceived resets the idle timer while responding.
func (d *Detector) OnTokenReceived() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.state != StateResponding {
		return
	}
	d.resetTimerLocked()
}

// OnResponseEnd transitions the detector to idle and clears the timer.
func (d *Detector) OnResponseEnd() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.state = StateIdle
	d.stopTimerLocked()
}

// State returns the detector's current state.
func (d *Detector) State() State {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.state
}

// Stop disables the detector and clears any active timer.
func (d *Detector) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.state = StateIdle
	d.stopTimerLocked()
}

func (d *Detector) resetTimerLocked() {
	d.stopTimerLocked()
	d.timer = time.AfterFunc(d.timeout, func() {
		d.mu.Lock()
		if d.state == StateResponding {
			d.state = StateStuck
			callback := d.onStuck
			d.mu.Unlock()
			if callback != nil {
				callback()
			}
			return
		}
		d.mu.Unlock()
	})
}

func (d *Detector) stopTimerLocked() {
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
