package port

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs-files"
	"github.com/multiverse-vcs/go-multiverse/core"
)

var (
	// ErrNoStart is returned when a starting commit is not found.
	ErrNoStart = errors.New("could not find starting commit")
)

type gitImporter struct {
	core *core.Core
	seen map[string]cid.Cid
}

// Static (compile time) check that gitImporter satisfies the Importer interface.
var _ Importer = (*gitImporter)(nil)

// NewGitImporter returns a new git importer.
func NewGitImporter(core *core.Core) *gitImporter {
	return &gitImporter{
		core: core,
		seen: map[string]cid.Cid{},
	}
}

// Import adds all commits from the repo at path.
func (im *gitImporter) Import(ctx context.Context, path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	opts := git.LogOptions{
		All:   true,
		Order: git.LogOrderCommitterTime,
	}

	iter, err := repo.Log(&opts)
	if err != nil {
		return err
	}

	var commits []plumbing.Hash
	cb := func(c *object.Commit) error {
		commits = append(commits, c.Hash)
		return nil
	}

	if err := iter.ForEach(cb); err != nil {
		return err
	}

	// add in reverse order so we don't recursively add parents
	for i := len(commits) - 1; i >= 0; i-- {
		c, err := repo.CommitObject(commits[i])
		if err != nil {
			return err
		}

		id, err := im.importGitCommit(ctx, c)
		if err != nil {
			return err
		}

		fmt.Printf("%s -> %s\n", c.Hash.String(), id.String())
	}

	return nil
}

// importGitCommit imports the git commit into an ipldmulti.Commit.
func (im *gitImporter) importGitCommit(ctx context.Context, c *object.Commit) (cid.Cid, error) {
	if id, ok := im.seen[c.Hash.String()]; ok {
		return id, nil
	}

	opts := core.CommitOptions{
		Message: c.Message,
	}

	cb := func(c *object.Commit) error {
		id, err := im.importGitCommit(ctx, c)
		if err != nil {
			return err
		}

		opts.Parents = append(opts.Parents, id)
		return nil
	}

	if err := c.Parents().ForEach(cb); err != nil {
		return cid.Cid{}, err
	}

	tree, err := c.Tree()
	if err != nil {
		return cid.Cid{}, nil
	}

	node, err := im.importGitTree(tree)
	if err != nil {
		return cid.Cid{}, nil
	}

	commit, err := im.core.Commit(ctx, node, &opts)
	if err != nil {
		return cid.Cid{}, nil
	}

	im.seen[c.Hash.String()] = commit.Cid()
	return commit.Cid(), nil
}

// importGitTree imports the git tree a files.Node.
func (im *gitImporter) importGitTree(tree *object.Tree) (files.Node, error) {
	dir := map[string]files.Node{}

	for _, entry := range tree.Entries {
		switch {
		case entry.Mode == filemode.Dir:
			other, err := tree.Tree(entry.Name)
			if err != nil {
				return nil, err
			}

			f, err := im.importGitTree(other)
			if err != nil {
				return nil, err
			}

			dir[entry.Name] = f
		case entry.Mode == filemode.Symlink:
			file, err := tree.File(entry.Name)
			if err != nil {
				return nil, err
			}

			r, err := file.Blob.Reader()
			if err != nil {
				return nil, err
			}

			target, err := ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}

			dir[entry.Name] = files.NewLinkFile(string(target), nil)
		case entry.Mode.IsFile():
			file, err := tree.File(entry.Name)
			if err != nil {
				return nil, err
			}

			r, err := file.Blob.Reader()
			if err != nil {
				return nil, err
			}

			dir[entry.Name] = files.NewReaderFile(r)
		}
	}

	return files.NewMapDirectory(dir), nil
}
