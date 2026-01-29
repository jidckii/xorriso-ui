package xorriso

import (
	"regexp"
	"strconv"
	"strings"
)

type Progress struct {
	Phase        string
	Percent      float64
	BytesWritten int64
	BytesTotal   int64
	Speed        string
	FIFOPercent  int
	ETA          string
}

// xorriso UPDATE pattern: "xorriso : UPDATE : 42.5% done"
var updatePercentRe = regexp.MustCompile(`(\d+\.?\d*)%\s+done`)
var updateFifoRe = regexp.MustCompile(`fifo\s+(\d+)%`)
var updateSpeedRe = regexp.MustCompile(`(\d+\.?\d*x[A-Z]+|\d+\.?\d*\s*[kMG]B/s)`)

// ParsePacifierLine tries to extract progress from an xorriso info line
func ParsePacifierLine(line string) (Progress, bool) {
	if !strings.Contains(line, "UPDATE") && !strings.Contains(line, "Writing:") && !strings.Contains(line, "Verifying") {
		return Progress{}, false
	}

	p := Progress{}

	// Determine phase
	switch {
	case strings.Contains(line, "Writing"):
		p.Phase = "writing"
	case strings.Contains(line, "Blanking"), strings.Contains(line, "Formatting"):
		p.Phase = "formatting"
	case strings.Contains(line, "Verifying"), strings.Contains(line, "check_media"):
		p.Phase = "verifying"
	default:
		p.Phase = "writing"
	}

	// Extract percent
	if m := updatePercentRe.FindStringSubmatch(line); m != nil {
		p.Percent, _ = strconv.ParseFloat(m[1], 64)
	}

	// Extract FIFO
	if m := updateFifoRe.FindStringSubmatch(line); m != nil {
		p.FIFOPercent, _ = strconv.Atoi(m[1])
	}

	// Extract speed
	if m := updateSpeedRe.FindStringSubmatch(line); m != nil {
		p.Speed = m[1]
	}

	return p, true
}
