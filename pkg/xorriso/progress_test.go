package xorriso

import (
	"math"
	"testing"
)

func TestParsePacifierLine(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantOK      bool
		wantPhase   string
		wantPercent float64
		wantSpeed   string
		wantFIFO    int
		wantETA     string
	}{
		{
			name:        "Writing с процентом, скоростью и fifo",
			line:        "xorriso : UPDATE : Writing: 42.5% done, 4.0xBD fifo 95%",
			wantOK:      true,
			wantPhase:   "writing",
			wantPercent: 42.5,
			wantSpeed:   "4.0xBD",
			wantFIFO:    95,
		},
		{
			name:        "Blanking → phase formatting",
			line:        "xorriso : UPDATE : Blanking: 10.0% done",
			wantOK:      true,
			wantPhase:   "formatting",
			wantPercent: 10.0,
		},
		{
			name:        "Verifying / check_media",
			line:        "xorriso : UPDATE : Verifying check_media: 78.2% done, 5200 kB/s",
			wantOK:      true,
			wantPhase:   "verifying",
			wantPercent: 78.2,
			wantSpeed:   "5200 kB/s",
		},
		{
			name:        "FIFO percent извлечение",
			line:        "xorriso : UPDATE : Writing: 0.1% done fifo 100%",
			wantOK:      true,
			wantPhase:   "writing",
			wantPercent: 0.1,
			wantFIFO:    100,
		},
		{
			name:        "Скорость 4.0xBD",
			line:        "xorriso : UPDATE : Writing: 50.0% done, 4.0xBD",
			wantOK:      true,
			wantPhase:   "writing",
			wantPercent: 50.0,
			wantSpeed:   "4.0xBD",
		},
		{
			name:        "Скорость kB/s",
			line:        "xorriso : UPDATE : Writing: 33.3% done, 5200 kB/s fifo 80%",
			wantOK:      true,
			wantPhase:   "writing",
			wantPercent: 33.3,
			wantSpeed:   "5200 kB/s",
			wantFIFO:    80,
		},
		{
			name:        "ETA remaining",
			line:        "xorriso : UPDATE : Writing: 60.0% done remaining 0:05:30",
			wantOK:      true,
			wantPhase:   "writing",
			wantPercent: 60.0,
			wantETA:     "0:05:30",
		},
		{
			name:   "Обычная строка — не прогресс",
			line:   "xorriso : NOTE : some informational message",
			wantOK: false,
		},
		{
			name:   "Пустая строка",
			line:   "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, ok := ParsePacifierLine(tt.line)
			if ok != tt.wantOK {
				t.Fatalf("ParsePacifierLine() ok = %v, хотели %v", ok, tt.wantOK)
			}

			if !tt.wantOK {
				return
			}

			if p.Phase != tt.wantPhase {
				t.Errorf("Phase = %q, хотели %q", p.Phase, tt.wantPhase)
			}

			if tt.wantPercent != 0 && math.Abs(p.Percent-tt.wantPercent) > 0.01 {
				t.Errorf("Percent = %f, хотели %f", p.Percent, tt.wantPercent)
			}

			if tt.wantSpeed != "" && p.Speed != tt.wantSpeed {
				t.Errorf("Speed = %q, хотели %q", p.Speed, tt.wantSpeed)
			}

			if tt.wantFIFO != 0 && p.FIFOPercent != tt.wantFIFO {
				t.Errorf("FIFOPercent = %d, хотели %d", p.FIFOPercent, tt.wantFIFO)
			}

			if tt.wantETA != "" && p.ETA != tt.wantETA {
				t.Errorf("ETA = %q, хотели %q", p.ETA, tt.wantETA)
			}
		})
	}
}
