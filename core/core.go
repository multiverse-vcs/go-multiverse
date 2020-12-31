// Package core implements common version control functions.
package core

import (
	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()
