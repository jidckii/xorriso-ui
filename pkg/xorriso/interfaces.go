package xorriso

import "context"

// Runner abstracts xorriso command execution for testability
type Runner interface {
	Run(ctx context.Context, args ...string) (*CmdResult, error)
	RunWithProgress(ctx context.Context, progressFn func(Progress), args ...string) (*CmdResult, error)
	Version(ctx context.Context) (string, error)
}
