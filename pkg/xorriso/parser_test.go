package xorriso

import (
	"reflect"
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

// --- ParsePktLine ---

func TestParsePktLine_ResultLine(t *testing.T) {
	pkt := ParsePktLine("R:0:some result text")
	if pkt == nil {
		t.Fatal("expected non-nil PktLine")
	}
	if pkt.Channel != 'R' {
		t.Errorf("expected Channel='R', got %q", pkt.Channel)
	}
	if pkt.Mode != 0 {
		t.Errorf("expected Mode=0, got %d", pkt.Mode)
	}
	if pkt.Text != "some result text" {
		t.Errorf("expected Text=%q, got %q", "some result text", pkt.Text)
	}
}

func TestParsePktLine_InfoLine(t *testing.T) {
	pkt := ParsePktLine("I:0:info text")
	if pkt == nil {
		t.Fatal("expected non-nil PktLine")
	}
	if pkt.Channel != 'I' {
		t.Errorf("expected Channel='I', got %q", pkt.Channel)
	}
	if pkt.Mode != 0 {
		t.Errorf("expected Mode=0, got %d", pkt.Mode)
	}
	if pkt.Text != "info text" {
		t.Errorf("expected Text=%q, got %q", "info text", pkt.Text)
	}
}

func TestParsePktLine_MarkLine(t *testing.T) {
	pkt := ParsePktLine("M:0:mark text")
	if pkt == nil {
		t.Fatal("expected non-nil PktLine")
	}
	if pkt.Channel != 'M' {
		t.Errorf("expected Channel='M', got %q", pkt.Channel)
	}
	if pkt.Text != "mark text" {
		t.Errorf("expected Text=%q, got %q", "mark text", pkt.Text)
	}
}

func TestParsePktLine_TooShort(t *testing.T) {
	pkt := ParsePktLine("R:0")
	if pkt != nil {
		t.Errorf("expected nil for too short line, got %+v", pkt)
	}
}

func TestParsePktLine_InvalidChannel(t *testing.T) {
	pkt := ParsePktLine("X:0:text")
	if pkt != nil {
		t.Errorf("expected nil for invalid channel 'X', got %+v", pkt)
	}
}

func TestParsePktLine_Mode1(t *testing.T) {
	// Mode 1 не должен обрезать trailing \n
	pkt := ParsePktLine("R:1:text\n")
	if pkt == nil {
		t.Fatal("expected non-nil PktLine")
	}
	if pkt.Mode != 1 {
		t.Errorf("expected Mode=1, got %d", pkt.Mode)
	}
	if pkt.Text != "text\n" {
		t.Errorf("expected Text=%q (with trailing newline), got %q", "text\n", pkt.Text)
	}
}

// --- ParsePktOutput ---

func TestParsePktOutput(t *testing.T) {
	output := "R:0:result line 1\nI:0:info line 1\nM:0:mark line 1\nR:0:result line 2\nI:0:info line 2\n"

	result := ParsePktOutput(output)

	expectedResult := []string{"result line 1", "result line 2"}
	expectedInfo := []string{"info line 1", "info line 2"}
	expectedMark := []string{"mark line 1"}

	if !reflect.DeepEqual(result.ResultLines, expectedResult) {
		t.Errorf("ResultLines: expected %v, got %v", expectedResult, result.ResultLines)
	}
	if !reflect.DeepEqual(result.InfoLines, expectedInfo) {
		t.Errorf("InfoLines: expected %v, got %v", expectedInfo, result.InfoLines)
	}
	if !reflect.DeepEqual(result.MarkLines, expectedMark) {
		t.Errorf("MarkLines: expected %v, got %v", expectedMark, result.MarkLines)
	}
}

// --- ParseDevices ---

func TestParseDevices(t *testing.T) {
	lines := []string{
		"0  -dev '/dev/sr0' rwrw-- : 'HL-DT-ST' 'BD-RE  WH16NS60'",
		"1  -dev '/dev/sr1' rwrw-- : 'ASUS' 'BW-16D1HT'",
		"some garbage line",
	}

	devices := ParseDevices(lines)

	if len(devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(devices))
	}

	if devices[0].Path != "/dev/sr0" {
		t.Errorf("device 0: expected Path=%q, got %q", "/dev/sr0", devices[0].Path)
	}
	if devices[0].Vendor != "HL-DT-ST" {
		t.Errorf("device 0: expected Vendor=%q, got %q", "HL-DT-ST", devices[0].Vendor)
	}
	if devices[0].Model != "BD-RE  WH16NS60" {
		t.Errorf("device 0: expected Model=%q, got %q", "BD-RE  WH16NS60", devices[0].Model)
	}

	if devices[1].Path != "/dev/sr1" {
		t.Errorf("device 1: expected Path=%q, got %q", "/dev/sr1", devices[1].Path)
	}
	if devices[1].Vendor != "ASUS" {
		t.Errorf("device 1: expected Vendor=%q, got %q", "ASUS", devices[1].Vendor)
	}
	if devices[1].Model != "BW-16D1HT" {
		t.Errorf("device 1: expected Model=%q, got %q", "BW-16D1HT", devices[1].Model)
	}
}

func TestParseDevices_Empty(t *testing.T) {
	devices := ParseDevices([]string{})
	if len(devices) != 0 {
		t.Errorf("expected 0 devices, got %d", len(devices))
	}

	devices = ParseDevices(nil)
	if len(devices) != 0 {
		t.Errorf("expected 0 devices for nil input, got %d", len(devices))
	}
}

// --- ParseSpeeds ---

func TestParseSpeeds(t *testing.T) {
	lines := []string{
		"Write speed  :   4234kB/s  (BD  1x)",
		"Write speed  :   8468kB/s  (BD  2x)",
		"some other line",
	}

	speeds := ParseSpeeds(lines)
	if len(speeds) != 2 {
		t.Fatalf("expected 2 speeds, got %d", len(speeds))
	}

	if speeds[0].WriteSpeed != 4234 {
		t.Errorf("speed 0: expected WriteSpeed=4234, got %f", speeds[0].WriteSpeed)
	}
	if speeds[0].DisplayName != "BD  1x" {
		t.Errorf("speed 0: expected DisplayName=%q, got %q", "BD  1x", speeds[0].DisplayName)
	}

	if speeds[1].WriteSpeed != 8468 {
		t.Errorf("speed 1: expected WriteSpeed=8468, got %f", speeds[1].WriteSpeed)
	}
	if speeds[1].DisplayName != "BD  2x" {
		t.Errorf("speed 1: expected DisplayName=%q, got %q", "BD  2x", speeds[1].DisplayName)
	}
}

func TestParseSpeeds_NoMatch(t *testing.T) {
	lines := []string{
		"some random line",
		"another line without speed info",
	}

	speeds := ParseSpeeds(lines)
	if len(speeds) != 0 {
		t.Errorf("expected 0 speeds, got %d", len(speeds))
	}
}

// --- ParseProfiles ---

func TestParseProfiles(t *testing.T) {
	lines := []string{
		"Profile      : 0x0043 (BD-RE)",
		"Profile      : 0x0041 (BD-R sequential recording) (current)",
		"Profile      : 0x0012 (DVD-RAM)",
		"some garbage",
	}

	profiles := ParseProfiles(lines)
	if len(profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %d", len(profiles))
	}

	if profiles[0].Name != "BD-RE" || profiles[0].Current {
		t.Errorf("profile 0: expected Name=%q Current=false, got Name=%q Current=%v",
			"BD-RE", profiles[0].Name, profiles[0].Current)
	}
	if profiles[1].Name != "BD-R sequential recording" || !profiles[1].Current {
		t.Errorf("profile 1: expected Name=%q Current=true, got Name=%q Current=%v",
			"BD-R sequential recording", profiles[1].Name, profiles[1].Current)
	}
	if profiles[2].Name != "DVD-RAM" || profiles[2].Current {
		t.Errorf("profile 2: expected Name=%q Current=false, got Name=%q Current=%v",
			"DVD-RAM", profiles[2].Name, profiles[2].Current)
	}
}

// --- ParseMediaSpace ---

func TestParseMediaSpace(t *testing.T) {
	lines := []string{
		"Media space  : 359844s",
	}

	blocks, err := ParseMediaSpace(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blocks != 359844 {
		t.Errorf("expected 359844, got %d", blocks)
	}
}

func TestParseMediaSpace_NotFound(t *testing.T) {
	lines := []string{
		"some other output",
	}

	blocks, err := ParseMediaSpace(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blocks != 0 {
		t.Errorf("expected 0 when not found, got %d", blocks)
	}
}

// --- ParseMediaBlocks ---

func TestParseMediaBlocks(t *testing.T) {
	lines := []string{
		"Media blocks : 100 readable , 200 writable , 300 overall",
	}

	readable, writable, overall := ParseMediaBlocks(lines)
	if readable != 100 {
		t.Errorf("expected readable=100, got %d", readable)
	}
	if writable != 200 {
		t.Errorf("expected writable=200, got %d", writable)
	}
	if overall != 300 {
		t.Errorf("expected overall=300, got %d", overall)
	}
}

// --- ParseCheckMediaResult ---

func TestParseCheckMediaResult(t *testing.T) {
	lines := []string{
		"quality scan: 1000 sectors, 995 readable, 5 unreadable",
		"MD5 mismatches :  3",
	}

	readErrors, md5 := ParseCheckMediaResult(lines)
	if readErrors != 5 {
		t.Errorf("expected readErrors=5, got %d", readErrors)
	}
	if md5 != 3 {
		t.Errorf("expected md5Mismatches=3, got %d", md5)
	}
}

func TestParseCheckMediaResult_Clean(t *testing.T) {
	lines := []string{
		"quality scan: 1000 sectors, 1000 readable, 0 unreadable",
		"MD5 mismatches :  0",
	}

	readErrors, md5 := ParseCheckMediaResult(lines)
	if readErrors != 0 {
		t.Errorf("expected readErrors=0, got %d", readErrors)
	}
	if md5 != 0 {
		t.Errorf("expected md5Mismatches=0, got %d", md5)
	}
}

// --- ParseMediaSummary ---

func TestParseMediaSummary(t *testing.T) {
	lines := []string{
		"Media summary: 2 sessions, 150000s data, 209844s free",
	}

	n := ParseMediaSummary(lines)
	if n != 2 {
		t.Errorf("expected 2 sessions, got %d", n)
	}
}

func TestParseMediaSummary_Zero(t *testing.T) {
	lines := []string{
		"Media summary: 0 sessions, 0s data, 359844s free",
	}

	n := ParseMediaSummary(lines)
	if n != 0 {
		t.Errorf("expected 0 sessions, got %d", n)
	}
}
