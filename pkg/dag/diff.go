package dag

import (
	"context"

	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
)

// ChangeType is used to specify the change made.
type ChangeType dagutils.ChangeType

const (
	Add    = ChangeType(dagutils.Add)
	Remove = ChangeType(dagutils.Remove)
	Mod    = ChangeType(dagutils.Mod)
)

// Diff returns a map of changes between two nodes.
func Diff(ctx context.Context, dag ipld.DAGService, old, new ipld.Node) (map[string]ChangeType, error) {
	changes, err := dagutils.Diff(ctx, dag, old, new)
	if err != nil {
		return nil, err
	}

	diffs := make(map[string]ChangeType)
	for _, change := range changes {
		if _, ok := diffs[change.Path]; ok {
			diffs[change.Path] = Mod
		} else if change.Path != "" {
			diffs[change.Path] = ChangeType(change.Type)
		}
	}

	return diffs, nil
}
