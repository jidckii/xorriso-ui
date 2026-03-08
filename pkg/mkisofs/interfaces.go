package mkisofs

import "context"

// ISOBuilder abstracts ISO image creation for testability
type ISOBuilder interface {
	BuildISO(ctx context.Context, opts BuildOpts, progressFn ProgressFn) error
}
