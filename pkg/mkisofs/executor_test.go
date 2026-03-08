package mkisofs

import (
	"bufio"
	"slices"
	"strings"
	"testing"

	"xorriso-ui/pkg/models"
)

func TestBuildArgs_AllFlags(t *testing.T) {
	e := NewExecutor("/usr/bin/mkisofs")
	opts := BuildOpts{
		UDF:        true,
		RockRidge:  true,
		Joliet:     true,
		HFSPlus:    true,
		Zisofs:     true,
		ISOLevel:   3,
		VolumeID:   "TEST",
		OutputPath: "/tmp/test.iso",
		Files: []FileMapping{
			{Source: "/home/user/file.txt", Dest: "/file.txt"},
		},
	}

	args := e.buildArgs(opts)

	expectedFlags := map[string]bool{
		"-udf":          false,
		"-r":            false,
		"-J":            false,
		"-hfsplus":      false,
		"-z":            false,
		"-iso-level":    false,
		"3":             false,
		"-V":            false,
		"TEST":          false,
		"-o":            false,
		"/tmp/test.iso": false,
		"-graft-points": false,
	}

	for _, arg := range args {
		if _, ok := expectedFlags[arg]; ok {
			expectedFlags[arg] = true
		}
	}

	for flag, found := range expectedFlags {
		if !found {
			t.Errorf("expected flag %q not found in args: %v", flag, args)
		}
	}
}

func TestBuildArgs_MinimalFlags(t *testing.T) {
	e := NewExecutor("/usr/bin/mkisofs")
	opts := BuildOpts{
		OutputPath: "/tmp/minimal.iso",
	}

	args := e.buildArgs(opts)

	// Не должно быть флагов файловых систем
	disallowed := []string{"-udf", "-r", "-J", "-hfsplus", "-z", "-iso-level", "-V"}
	for _, d := range disallowed {
		for _, arg := range args {
			if arg == d {
				t.Errorf("unexpected flag %q in minimal args: %v", d, args)
			}
		}
	}

	// Обязательные флаги должны быть
	mustHave := []string{"-o", "/tmp/minimal.iso", "-graft-points"}
	for _, m := range mustHave {
		found := slices.Contains(args, m)
		if !found {
			t.Errorf("expected flag %q not found in args: %v", m, args)
		}
	}
}

func TestBuildArgs_FileMappings(t *testing.T) {
	e := NewExecutor("/usr/bin/mkisofs")
	opts := BuildOpts{
		OutputPath: "/tmp/test.iso",
		Files: []FileMapping{
			{Source: "/home/user/doc.pdf", Dest: "/docs/doc.pdf"},
			{Source: "/home/user/pic.jpg", Dest: "/images/pic.jpg"},
		},
	}

	args := e.buildArgs(opts)

	// Проверяем graft-points формат: dest=source
	expected := map[string]bool{
		"/docs/doc.pdf=/home/user/doc.pdf":   false,
		"/images/pic.jpg=/home/user/pic.jpg": false,
	}

	for _, arg := range args {
		if _, ok := expected[arg]; ok {
			expected[arg] = true
		}
	}

	for mapping, found := range expected {
		if !found {
			t.Errorf("expected graft-point %q not found in args: %v", mapping, args)
		}
	}
}

func TestFileMappingsFromEntries(t *testing.T) {
	entries := []models.FileEntry{
		{SourcePath: "/home/user/file1.txt", DestPath: "/file1.txt", Name: "file1.txt"},
		{SourcePath: "/home/user/dir/file2.txt", DestPath: "/dir/file2.txt", Name: "file2.txt"},
	}

	mappings := FileMappingsFromEntries(entries)

	if len(mappings) != 2 {
		t.Fatalf("expected 2 mappings, got %d", len(mappings))
	}

	if mappings[0].Source != "/home/user/file1.txt" {
		t.Errorf("mapping 0: expected Source=%q, got %q", "/home/user/file1.txt", mappings[0].Source)
	}
	if mappings[0].Dest != "/file1.txt" {
		t.Errorf("mapping 0: expected Dest=%q, got %q", "/file1.txt", mappings[0].Dest)
	}

	if mappings[1].Source != "/home/user/dir/file2.txt" {
		t.Errorf("mapping 1: expected Source=%q, got %q", "/home/user/dir/file2.txt", mappings[1].Source)
	}
	if mappings[1].Dest != "/dir/file2.txt" {
		t.Errorf("mapping 1: expected Dest=%q, got %q", "/dir/file2.txt", mappings[1].Dest)
	}
}

func TestScanMkisofsLines(t *testing.T) {
	// mkisofs использует \r для обновления прогресса на месте, \n для обычных строк
	input := " 10.02% done\r 20.05% done\rTotal bytes: 1234\n"

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(scanMkisofsLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}

	expected := []string{
		" 10.02% done",
		" 20.05% done",
		"Total bytes: 1234",
	}

	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d: %v", len(expected), len(lines), lines)
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}
