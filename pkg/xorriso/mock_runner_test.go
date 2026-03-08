package xorriso

import "context"

// MockRunner is a test double for xorriso.Runner
type MockRunner struct {
	RunFn             func(ctx context.Context, args ...string) (*CmdResult, error)
	RunWithProgressFn func(ctx context.Context, progressFn func(Progress), args ...string) (*CmdResult, error)
	VersionFn         func(ctx context.Context) (string, error)
}

func (m *MockRunner) Run(ctx context.Context, args ...string) (*CmdResult, error) {
	if m.RunFn != nil {
		return m.RunFn(ctx, args...)
	}
	return &CmdResult{}, nil
}

func (m *MockRunner) RunWithProgress(ctx context.Context, progressFn func(Progress), args ...string) (*CmdResult, error) {
	if m.RunWithProgressFn != nil {
		return m.RunWithProgressFn(ctx, progressFn, args...)
	}
	return &CmdResult{}, nil
}

func (m *MockRunner) Version(ctx context.Context) (string, error) {
	if m.VersionFn != nil {
		return m.VersionFn(ctx)
	}
	return "xorriso 1.5.6", nil
}
