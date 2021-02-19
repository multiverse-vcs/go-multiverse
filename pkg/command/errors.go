package command

import (
	"errors"
)

// ErrDialRPC is an error message for failed RPC connections.
var ErrDialRPC = errors.New(`
Could not connect to local RPC server.
Make sure the Multiverse daemon is up.
See 'multi help daemon' for more info.
`)
