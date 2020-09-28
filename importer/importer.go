// Package importer implements methods for importing repos
// from other VCS systems into Multiverse.
package importer

// Importer defines an interface for importers.
type Importer interface {
	Import(path string) error
}