package mkisofs

import "context"

// MockISOBuilder is a test double for mkisofs.ISOBuilder
type MockISOBuilder struct {
	BuildISOFn func(ctx context.Context, opts BuildOpts, progressFn ProgressFn) error
}

func (m *MockISOBuilder) BuildISO(ctx context.Context, opts BuildOpts, progressFn ProgressFn) error {
	if m.BuildISOFn != nil {
		return m.BuildISOFn(ctx, opts, progressFn)
	}
	return nil
}
