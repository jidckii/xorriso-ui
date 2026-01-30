package xorriso

import (
	"regexp"
	"strconv"
	"strings"

	"xorriso-ui/pkg/models"
)

type PktLine struct {
	Channel byte
	Mode    int
	Text    string
}

// ParsePktLine parses a single pkt_output line like "R:0:text"
func ParsePktLine(line string) *PktLine {
	if len(line) < 4 {
		return nil
	}
	if line[1] != ':' || line[3] != ':' {
		return nil
	}

	ch := line[0]
	if ch != 'R' && ch != 'I' && ch != 'M' {
		return nil
	}

	mode, err := strconv.Atoi(string(line[2]))
	if err != nil {
		return nil
	}

	text := line[4:]
	// Mode 0: trailing \n is extra (strip it)
	if mode == 0 {
		text = strings.TrimRight(text, "\n")
	}

	return &PktLine{Channel: ch, Mode: mode, Text: text}
}

// ParsePktOutput parses all pkt_output lines into a CmdResult
func ParsePktOutput(output string) *CmdResult {
	result := &CmdResult{}
	for _, line := range strings.Split(output, "\n") {
		pkt := ParsePktLine(line)
		if pkt == nil {
			continue
		}
		switch pkt.Channel {
		case 'R':
			result.ResultLines = append(result.ResultLines, pkt.Text)
		case 'I':
			result.InfoLines = append(result.InfoLines, pkt.Text)
		case 'M':
			result.MarkLines = append(result.MarkLines, pkt.Text)
		}
	}
	return result
}

// Device line pattern: N  -dev 'PATH' rwrw-- : 'VENDOR' 'MODEL'
var deviceLineRe = regexp.MustCompile(`^\s*(\d+)\s+-dev\s+'([^']+)'\s+\S+\s+:\s+'([^']*)'\s+'([^']*)'`)

// ParseDevices parses output of -devices or -device_links
func ParseDevices(lines []string) []models.Device {
	var devices []models.Device
	for _, line := range lines {
		matches := deviceLineRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		devices = append(devices, models.Device{
			Path:   matches[2],
			Vendor: strings.TrimSpace(matches[3]),
			Model:  strings.TrimSpace(matches[4]),
		})
	}
	return devices
}

// ParseSpeeds parses output of -list_speeds
// Lines like: "Write speed  :   4234kB/s  (BD  1x)"
var speedLineRe = regexp.MustCompile(`(\d+)kB/s\s+\(([^)]+)\)`)

func ParseSpeeds(lines []string) []models.SpeedDescriptor {
	var speeds []models.SpeedDescriptor
	for _, line := range lines {
		if !strings.Contains(line, "kB/s") {
			continue
		}
		matches := speedLineRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		kbps, _ := strconv.ParseFloat(matches[1], 64)
		speeds = append(speeds, models.SpeedDescriptor{
			WriteSpeed:  kbps,
			DisplayName: strings.TrimSpace(matches[2]),
		})
	}
	return speeds
}

// ParseProfiles parses output of -list_profiles
// Actual xorriso pkt_output format: "Profile      : 0x0043 (BD-RE)"
// With current profile marked: "Profile      : 0x0041 (BD-R sequential recording) (current)"
var profileLineRe = regexp.MustCompile(`Profile\s+:\s+0x([0-9A-Fa-f]+)\s+\(([^)]+)\)`)

func ParseProfiles(lines []string) []models.MediaProfile {
	var profiles []models.MediaProfile
	for _, line := range lines {
		matches := profileLineRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		current := strings.Contains(line, "(current)")
		profiles = append(profiles, models.MediaProfile{
			Name:    strings.TrimSpace(matches[2]),
			Current: current,
		})
	}
	return profiles
}

// ParseMediaSpace parses output of -tell_media_space
// Line: "Media space  : NNNNN  (free blocks)"
var mediaSpaceRe = regexp.MustCompile(`(\d+)s\s+\(`)

func ParseMediaSpace(lines []string) (freeBlocks int64, err error) {
	for _, line := range lines {
		matches := mediaSpaceRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		freeBlocks, err = strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return 0, err
		}
		return freeBlocks, nil
	}
	return 0, nil
}
