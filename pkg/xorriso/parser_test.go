package xorriso

import (
	"testing"

	"xorriso-ui/pkg/models"
)

func TestParseTOCSessions(t *testing.T) {
	lines := []string{
		"TOC layout   : Idx ,  sbsector ,       Size , Volume Id",
		"ISO session  :   1 ,         0 ,    150000s , MY_DISC",
		"ISO session  :   2 ,    150000 ,     50000s , MY_DISC_2",
	}

	sessions := ParseTOCSessions(lines)

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}

	expected := []models.Session{
		{Number: 1, StartLBA: 0, Size: 150000, VolumeID: "MY_DISC"},
		{Number: 2, StartLBA: 150000, Size: 50000, VolumeID: "MY_DISC_2"},
	}

	for i, s := range sessions {
		if s.Number != expected[i].Number {
			t.Errorf("session %d: expected Number=%d, got %d", i, expected[i].Number, s.Number)
		}
		if s.StartLBA != expected[i].StartLBA {
			t.Errorf("session %d: expected StartLBA=%d, got %d", i, expected[i].StartLBA, s.StartLBA)
		}
		if s.Size != expected[i].Size {
			t.Errorf("session %d: expected Size=%d, got %d", i, expected[i].Size, s.Size)
		}
		if s.VolumeID != expected[i].VolumeID {
			t.Errorf("session %d: expected VolumeID=%q, got %q", i, expected[i].VolumeID, s.VolumeID)
		}
	}
}

func TestParseTOCSessions_Empty(t *testing.T) {
	lines := []string{
		"TOC layout   : Idx ,  sbsector ,       Size , Volume Id",
	}

	sessions := ParseTOCSessions(lines)
	if len(sessions) != 0 {
		t.Fatalf("expected 0 sessions, got %d", len(sessions))
	}
}

func TestParseTOCSessions_SingleSession(t *testing.T) {
	lines := []string{
		"ISO session  :   1 ,         0 ,    359844s , BACKUP_2026",
	}

	sessions := ParseTOCSessions(lines)
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].VolumeID != "BACKUP_2026" {
		t.Errorf("expected VolumeID=%q, got %q", "BACKUP_2026", sessions[0].VolumeID)
	}
}
