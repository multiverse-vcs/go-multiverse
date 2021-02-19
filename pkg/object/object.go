// Package object contains object definitions.
package object

import (
	"time"

	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/polydawn/refmt/obj/atlas"
)

// timeAtlasEntry allows encoding and decoding of time structs.
var timeAtlasEntry = atlas.BuildEntry(time.Time{}).
	Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(t time.Time) (string, error) {
			return t.Format(time.RFC3339), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
		func(t string) (time.Time, error) {
			return time.Parse(time.RFC3339, t)
		})).
	Complete()

func init() {
	cbornode.RegisterCborType(timeAtlasEntry)
	cbornode.RegisterCborType(Author{})
	cbornode.RegisterCborType(Commit{})
	cbornode.RegisterCborType(Repository{})
}
