package dag

import (
	"context"
	"errors"
	"io"

	cid "github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	car "github.com/ipld/go-car"
	"github.com/ipld/go-car/util"
)

type carWriter struct {
	ds ipld.NodeGetter
	w  io.Writer
}

func (cw *carWriter) getLinks(ctx context.Context, id cid.Cid) ([]*ipld.Link, error) {
	node, err := cw.ds.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := util.LdWrite(cw.w, id.Bytes(), node.RawData()); err != nil {
		return nil, err
	}

	return node.Links(), nil
}

// WriteCar writes the dag into the given writer starting at head and stopping at refs.
func WriteCar(ctx context.Context, ds ipld.NodeGetter, head cid.Cid, refs *cid.Set, w io.Writer) error {
	h := &car.CarHeader{
		Roots:   []cid.Cid{head},
		Version: 1,
	}

	if err := car.WriteHeader(h, w); err != nil {
		return err
	}

	cw := carWriter{
		ds: ds,
		w:  w,
	}

	return merkledag.Walk(ctx, cw.getLinks, head, refs.Visit)
}

// ReadCar reads the car into the given dag and returns the root cid.
func ReadCar(bs blockstore.Blockstore, r io.Reader) (cid.Cid, error) {
	cr, err := car.NewCarReader(r)
	if err != nil {
		return cid.Cid{}, err
	}

	if len(cr.Header.Roots) != 1 {
		return cid.Cid{}, errors.New("unexpected header roots")
	}

	// load blocks slowly or badger will return an error
	for {
		block, err := cr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return cid.Cid{}, err
		}

		if err := bs.Put(block); err != nil {
			return cid.Cid{}, err
		}
	}

	return cr.Header.Roots[0], nil
}
