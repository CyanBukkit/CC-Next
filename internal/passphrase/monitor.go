package passphrase

import (
	"strings"
)

// Monitor scans streaming text for a passphrase and filters it out.
type Monitor struct {
	phrase    string
	ringBuf   []byte // ring buffer for cross-chunk matching
	bufPos    int
	bufFilled bool
	found     bool
}

// NewMonitor creates a new Monitor for the supplied passphrase.
func NewMonitor(phrase string) *Monitor {
	bufSize := len(phrase) - 1
	if bufSize < 0 {
		bufSize = 0
	}
	return &Monitor{
		phrase:  phrase,
		ringBuf: make([]byte, bufSize),
	}
}

// Feed processes a text chunk. Returns the filtered text (passphrase removed)
// and whether the passphrase was found in this chunk.
func (m *Monitor) Feed(chunk string) (filtered string, found bool) {
	p := len(m.phrase)
	if p == 0 {
		return chunk, false
	}

	prefix := m.bufferString()
	combined := prefix + chunk

	hit := false
	if strings.Contains(combined, m.phrase) {
		hit = true
		m.found = true
		combined = strings.ReplaceAll(combined, m.phrase, "")
	}

	keep := p - 1
	if keep < 0 {
		keep = 0
	}
	outputLen := len(combined) - keep
	if outputLen < 0 {
		outputLen = 0
	}

	output := combined[:outputLen]
	suffix := combined[outputLen:]
	m.setBuffer(suffix)

	return output, hit
}

// Found returns true if the passphrase was ever detected.
func (m *Monitor) Found() bool {
	return m.found
}

// Flush drains any remaining buffered bytes and resets the buffer.
func (m *Monitor) Flush() string {
	remaining := m.bufferString()
	m.ringBuf = make([]byte, len(m.ringBuf))
	m.bufPos = 0
	m.bufFilled = false
	return remaining
}

// bufferString returns the current contents of the ring buffer in order.
func (m *Monitor) bufferString() string {
	if m.bufFilled {
		return string(m.ringBuf[m.bufPos:]) + string(m.ringBuf[:m.bufPos])
	}
	return string(m.ringBuf[:m.bufPos])
}

// setBuffer stores the given suffix in the ring buffer.
func (m *Monitor) setBuffer(s string) {
	if len(m.ringBuf) == 0 {
		m.bufPos = 0
		m.bufFilled = false
		return
	}

	if len(s) > len(m.ringBuf) {
		s = s[len(s)-len(m.ringBuf):]
	}

	for i := 0; i < len(s); i++ {
		m.ringBuf[m.bufPos] = s[i]
		m.bufPos++
		if m.bufPos >= len(m.ringBuf) {
			m.bufPos = 0
			m.bufFilled = true
		}
	}
}
