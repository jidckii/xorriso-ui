package xorriso

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type Executor struct {
	binaryPath string
	mu         sync.Mutex
}

func NewExecutor(binaryPath string) *Executor {
	return &Executor{binaryPath: binaryPath}
}

type CmdResult struct {
	ResultLines []string
	InfoLines   []string
	MarkLines   []string
	ExitCode    int
	RawOutput   string
}

// Run executes a short xorriso command and returns parsed result
func (e *Executor) Run(ctx context.Context, args ...string) (*CmdResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Prepend pkt_output for machine-readable output
	fullArgs := append([]string{"-pkt_output", "on"}, args...)
	cmd := exec.CommandContext(ctx, e.binaryPath, fullArgs...)

	output, err := cmd.CombinedOutput()
	rawOutput := string(output)

	result := ParsePktOutput(rawOutput)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("failed to execute xorriso: %w", err)
		}
	}
	result.RawOutput = rawOutput

	return result, nil
}

// RunWithProgress executes a long operation with real-time progress updates
func (e *Executor) RunWithProgress(ctx context.Context, progressFn func(Progress), args ...string) (*CmdResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	fullArgs := append([]string{"-pkt_output", "on"}, args...)
	cmd := exec.CommandContext(ctx, e.binaryPath, fullArgs...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	cmd.Stderr = cmd.Stdout // merge stderr into stdout

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start xorriso: %w", err)
	}

	result := &CmdResult{}
	var rawLines []string

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		rawLines = append(rawLines, line)

		pktLine := ParsePktLine(line)
		if pktLine == nil {
			continue
		}

		switch pktLine.Channel {
		case 'R':
			result.ResultLines = append(result.ResultLines, pktLine.Text)
		case 'I':
			result.InfoLines = append(result.InfoLines, pktLine.Text)
			// Check for progress updates
			if p, ok := ParsePacifierLine(pktLine.Text); ok && progressFn != nil {
				progressFn(p)
			}
		case 'M':
			result.MarkLines = append(result.MarkLines, pktLine.Text)
		}
	}

	result.RawOutput = strings.Join(rawLines, "\n")

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			return result, fmt.Errorf("xorriso process error: %w", err)
		}
	}

	return result, nil
}

// Version returns the xorriso version string
func (e *Executor) Version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, e.binaryPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get xorriso version: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
