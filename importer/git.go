package importer

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipld-format"
	"github.com/ipfs/interface-go-ipfs-core"
	"github.com/yondero/multiverse/commit"
)

// GitImporter implements the Importer interface.
type GitImporter struct {
	api  iface.CoreAPI
	repo *git.Repository
}

// NewGitImporter creates a new importer for the git repo in the given path.
func NewGitImporter(api iface.CoreAPI) *GitImporter {
	return &GitImporter{api: api}
}

// Import adds commits starting at the current repo head.
func (i *GitImporter) Import(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	i.repo = repo

	head, err := i.repo.Head()
	if err != nil {
		return err
	}

	_, err = i.addCommit(head.Hash())
	if err != nil {
		return err
	}

	return nil
}

func (i *GitImporter) addCommit(h plumbing.Hash) (format.Node, error) {
	gc, err := i.repo.CommitObject(h)
	if err != nil {
		return nil, err
	}

	var parents []cid.Cid
	for _, ph := range gc.ParentHashes {
		node, err := i.addCommit(ph)
		if err != nil {
			return nil, err
		}

		parents = append(parents, node.Cid())
	}

	tree, err := i.addTree(gc.TreeHash)
	if err != nil {
		return nil, err
	}

	p, err := i.api.Unixfs().Add(context.TODO(), tree)
	if err != nil {
		return nil, err
	}

	c := commit.Commit{Message: gc.Message, Parents: parents, Tree: p.Root()}

	node, err := c.Node()
	if err != nil {
		return nil, err
	}

	if err := i.api.Dag().Pinning().Add(context.TODO(), node); err != nil {
		return nil, err
	}

	fmt.Println(node.Cid().String())
	return node, nil
}

func (i *GitImporter) addTree(h plumbing.Hash) (files.Node, error) {
	t, err := i.repo.TreeObject(h)
	if err != nil {
		return nil, err
	}

	dir := make(map[string]files.Node)
	for _, e := range t.Entries {
		f, err := i.addFile(t, &e)
		if err != nil {
			return nil, err
		}

		dir[e.Name] = f
	}

	return files.NewMapDirectory(dir), nil
}

func (i *GitImporter) addFile(t *object.Tree, e *object.TreeEntry) (files.Node, error) {
	if !e.Mode.IsFile() {
		return i.addTree(e.Hash)
	}

	f, err := t.TreeEntryFile(e)
	if err != nil {
		return nil, err
	}

	r, err := f.Reader()
	if err != nil {
		return nil, err
	}

	return files.NewReaderFile(r), nil
}