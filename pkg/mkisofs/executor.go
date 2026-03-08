package mkisofs

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"xorriso-ui/pkg/models"
)

// Executor manages mkisofs subprocess execution
type Executor struct {
	binaryPath string
}

// NewExecutor creates a new mkisofs executor
func NewExecutor(binaryPath string) *Executor {
	return &Executor{binaryPath: binaryPath}
}

// BuildOpts contains options for building an ISO image
type BuildOpts struct {
	OutputPath string
	VolumeID   string
	UDF        bool
	RockRidge  bool
	Joliet     bool
	HFSPlus    bool
	Zisofs     bool
	ISOLevel   int
	Files      []FileMapping
}

// FileMapping maps a source path to a destination path in the ISO
type FileMapping struct {
	Source string
	Dest   string
}

// ProgressFn callback for progress updates
type ProgressFn func(percent float64)

// progressRegex matches mkisofs progress lines like " 10.02% done"
var progressRegex = regexp.MustCompile(`(\d+\.?\d*)%\s+done`)

// BuildISO creates an ISO image with the given options.
// Returns error on failure. Progress is reported via progressFn.
func (e *Executor) BuildISO(ctx context.Context, opts BuildOpts, progressFn ProgressFn) error {
	args := e.buildArgs(opts)

	cmd := exec.CommandContext(ctx, e.binaryPath, args...)

	// mkisofs writes progress to stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start mkisofs: %w", err)
	}

	// Parse stderr for progress
	scanner := bufio.NewScanner(stderr)
	scanner.Split(scanMkisofsLines)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := progressRegex.FindStringSubmatch(line); len(matches) > 1 {
			if pct, err := strconv.ParseFloat(matches[1], 64); err == nil && progressFn != nil {
				progressFn(pct)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("mkisofs failed: %w", err)
	}

	return nil
}

// buildArgs constructs mkisofs command-line arguments
func (e *Executor) buildArgs(opts BuildOpts) []string {
	var args []string

	if opts.UDF {
		args = append(args, "-udf")
	}
	if opts.RockRidge {
		args = append(args, "-r")
	}
	if opts.Joliet {
		args = append(args, "-J")
	}
	if opts.HFSPlus {
		args = append(args, "-hfsplus")
	}
	if opts.Zisofs {
		args = append(args, "-z")
	}
	if opts.ISOLevel > 0 {
		args = append(args, "-iso-level", strconv.Itoa(opts.ISOLevel))
	}
	if opts.VolumeID != "" {
		args = append(args, "-V", opts.VolumeID)
	}

	args = append(args, "-o", opts.OutputPath)
	args = append(args, "-graft-points")

	for _, f := range opts.Files {
		args = append(args, f.Dest+"="+f.Source)
	}

	return args
}

// FileMappingsFromEntries converts project entries to file mappings
func FileMappingsFromEntries(entries []models.FileEntry) []FileMapping {
	mappings := make([]FileMapping, 0, len(entries))
	for _, e := range entries {
		mappings = append(mappings, FileMapping{
			Source: e.SourcePath,
			Dest:   e.DestPath,
		})
	}
	return mappings
}

// scanMkisofsLines is a custom split function for bufio.Scanner
// that handles \r-delimited progress lines from mkisofs
func scanMkisofsLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// Find \r or \n
	for i := range data {
		if data[i] == '\n' || data[i] == '\r' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
