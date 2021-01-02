package html

import (
	"context"
	"html/template"
	"io"

	"github.com/ipfs/go-path"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/node"
)

var blobView = template.Must(template.ParseFiles("html/index.html", "html/blob.html"))

type blobModel struct {
	Data string
	Path string
	Repo string
	Util *util
}

// Blob renders the blob page.
func Blob(ctx context.Context, w io.Writer, node *node.Node, repo, file string) error {
	id, err := node.Repo().Get(repo)
	if err != nil {
		return err
	}

	p, err := path.FromSegments("/ipfs/", id.String(), "tree", file)
	if err != nil {
		return err
	}

	blob, err := node.ResolvePath(ctx, p)
	if err != nil {
		return err
	}

	data, err := core.Cat(ctx, node, blob.Cid())
	if err != nil {
		return err
	}

	model := blobModel{
		Data: data,
		Path: file,
		Repo: repo,
		Util: &util{node},
	}

	return blobView.Execute(w, &model)
}
