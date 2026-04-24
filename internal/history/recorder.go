package history

import (
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Recorder wraps a History and converts monitor.Change values into entries.
type Recorder struct {
	h *History
}

// NewRecorder returns a Recorder backed by h.
func NewRecorder(h *History) *Recorder {
	return &Recorder{h: h}
}

// Record persists all changes from a monitor scan result.
func (r *Recorder) Record(changes []monitor.Change) error {
	for _, c := range changes {
		action := "opened"
		if c.Closed {
			action = "closed"
		}
		e := Entry{
			Timestamp: time.Now().UTC(),
			Port:      c.Port.Port,
			Proto:     c.Port.Proto,
			Action:    action,
		}
		if err := r.h.Add(e); err != nil {
			return err
		}
	}
	return nil
}
