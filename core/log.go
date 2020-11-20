package core

import (
	"fmt"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/object"
	"github.com/multiverse-vcs/go-multiverse/util"
)

// LogDateFormat is the format used when logging the commit.
const LogDateFormat = "Mon Jan 2 15:04:05 2006 -0700"

// Log prints commit history starting at the current head.
func (c *Context) Log(w io.Writer) error {
	cb := func(id cid.Cid, commit *object.Commit) bool {
		fmt.Fprintf(w, "%scommit %s", util.ColorYellow, id.String())

		if id == c.config.Head {
			fmt.Fprintf(w, " (%sHEAD%s)", util.ColorRed, util.ColorYellow)
		}

		if id == c.config.Base {
			fmt.Fprintf(w, " (%sBASE%s)", util.ColorGreen, util.ColorYellow)
		}

		fmt.Fprintf(w, "%s\n", util.ColorReset)
		fmt.Fprintf(w, "Date: %s\n", commit.Date.Format(LogDateFormat))
		fmt.Fprintf(w, "\n\t%s\n\n", commit.Message)
		return true
	}

	if _, err := c.Walk(c.config.Head, cb); err != nil {
		return err
	}

	return nil
}
