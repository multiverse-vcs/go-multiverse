// Package object contains Multiverse object definitions.
package object

import (
	"time"

	"github.com/ipfs/go-ipld-cbor"
	"github.com/polydawn/refmt/obj/atlas"
)

// timeAtlasEntry allows encoding and decoding of time structs.
var timeAtlasEntry = atlas.BuildEntry(time.Time{}).
	Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(t time.Time) ([]byte, error) {
			return t.MarshalText()
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
		func(data []byte) (time.Time, error) {
			return time.Parse(time.RFC3339, string(data))
		})).
	Complete()

// register all types here so they can be encoded and decoded
func init() {
	cbornode.RegisterCborType(timeAtlasEntry)
	cbornode.RegisterCborType(Commit{})
}
