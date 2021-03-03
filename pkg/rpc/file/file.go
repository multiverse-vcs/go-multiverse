package file

import (
	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
)

// Service wraps a remote and provides RPC.
type Service struct {
	*remote.Server
}
