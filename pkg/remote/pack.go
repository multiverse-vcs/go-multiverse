package remote

import (
	"bytes"
	"context"
	"errors"

	cid "github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	car "github.com/ipld/go-car"
)

// BuildPack returns a pack containing nodes missing from old and not in have.
func BuildPack(ctx context.Context, dag ipld.DAGService, heads *cid.Set, old, new cid.Cid) ([]byte, error) {
	if err := verifyPack(ctx, dag, old, new); err != nil {
		return nil, err
	}

	root := []cid.Cid{new}
	walk := func(node ipld.Node) ([]*ipld.Link, error) {
		if heads.Has(node.Cid()) {
			return nil, nil
		}

		return node.Links(), nil
	}

	var pack bytes.Buffer
	if err := car.WriteCarWithWalker(ctx, dag, root, &pack, walk); err != nil {
		return nil, err
	}

	return pack.Bytes(), nil
}

// LoadPack adds the pack data to the dag and verifies the root is valid.
func LoadPack(ctx context.Context, dag ipld.DAGService, bs blockstore.Blockstore, data []byte, old cid.Cid) (cid.Cid, error) {
	head, err := car.LoadCar(bs, bytes.NewReader(data))
	if err != nil {
		return cid.Cid{}, err
	}

	if len(head.Roots) != 1 {
		return cid.Cid{}, errors.New("unexpected pack roots")
	}

	if err := verifyPack(ctx, dag, old, head.Roots[0]); err != nil {
		return cid.Cid{}, err
	}

	return head.Roots[0], nil
}

// verifyPack checks if the new branch head contains the old branch head.
func verifyPack(ctx context.Context, dag ipld.DAGService, old, new cid.Cid) error {
	if !old.Defined() {
		return nil
	}

	refs := cid.NewSet()
	walk := merkledag.GetLinksWithDAG(dag)

	if err := merkledag.Walk(ctx, walk, new, refs.Visit); err != nil {
		return err
	}

	if !refs.Has(old) {
		return errors.New("remote is ahead of local")
	}

	return nil
}
