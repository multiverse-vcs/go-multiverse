// Package port implements methods to import and export repos.
package port

import (
	"context"
)

// Importer is an interface for importing repos.
type Importer interface {
	// Import adds all commits from the repo at path.
	Import(ctx context.Context, path string) error
}
