// Package config contains configuration definitions.
package config

import (
	"github.com/ipfs/go-cid"
)

// DefaultBranch is the name of the default branch.
const DefaultBranch = "default"

// Config contains local repo info.
type Config struct {
	// Base is the commit changes are based on.
	Base cid.Cid `json:"base"`
	// Head is the tip of the current branch.
	Head cid.Cid `json:"head"`
	// Branch is the name of the current branch.
	Branch string `json:"branch"`
}
