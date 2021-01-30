// Package git contains methods for importing Git repositories.
package git

import (
	"context"
	"io/ioutil"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	ufsio "github.com/ipfs/go-unixfs/io"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

// DateFormat is the format to store dates in.
const DateFormat = "Mon Jan 02 15:04:05 2006 -0700"

// importer adds objects from a git repo to a dag.
type importer struct {
	ctx      context.Context
	dag      ipld.DAGService
	name     string
	repo     *git.Repository
	branches map[string]cid.Cid
	tags     map[string]cid.Cid
}

// ImportFromFS is a helper to import a git repo from a directory.
func ImportFromFS(ctx context.Context, dag ipld.DAGService, name, dir string) (cid.Cid, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return cid.Cid{}, err
	}

	importer := NewImporter(ctx, dag, repo, name)
	return importer.AddRepository()
}

// ImportFromURL is a helper to import a git repo from a url.
func ImportFromURL(ctx context.Context, dag ipld.DAGService, name, url string) (cid.Cid, error) {
	opts := git.CloneOptions{
		URL: url,
	}

	repo, err := git.Clone(memory.NewStorage(), nil, &opts)
	if err != nil {
		return cid.Cid{}, err
	}

	importer := NewImporter(ctx, dag, repo, name)
	return importer.AddRepository()
}

// NewImporter returns an importer for the given repo.
func NewImporter(ctx context.Context, dag ipld.DAGService, repo *git.Repository, name string) *importer {
	return &importer{
		ctx:      ctx,
		dag:      dag,
		name:     name,
		repo:     repo,
		branches: make(map[string]cid.Cid),
		tags:     make(map[string]cid.Cid),
	}
}

// AddRepository adds all branches and tags to the dag.
func (i *importer) AddRepository() (cid.Cid, error) {
	tags, err := i.repo.Tags()
	if err != nil {
		return cid.Cid{}, err
	}

	branches, err := i.repo.Branches()
	if err != nil {
		return cid.Cid{}, err
	}

	if err := branches.ForEach(i.AddBranch); err != nil {
		return cid.Cid{}, err
	}

	if err := tags.ForEach(i.AddTag); err != nil {
		return cid.Cid{}, err
	}

	mrepo := data.NewRepository(i.name)
	mrepo.Branches = i.branches
	mrepo.Tags = i.tags

	return data.AddRepository(i.ctx, i.dag, mrepo)
}

// AddBranch adds the branch with the given ref to the dag.
func (i *importer) AddBranch(ref *plumbing.Reference) error {
	id, err := i.AddCommit(ref.Hash())
	if err != nil {
		return err
	}

	i.branches[string(ref.Name())] = id
	return nil
}

// AddTag adds the tag with the given ref to the dag.
func (i *importer) AddTag(ref *plumbing.Reference) error {
	id, err := i.AddCommit(ref.Hash())
	if err != nil {
		return err
	}

	i.tags[string(ref.Name())] = id
	return nil
}

// AddCommit adds the commit with the given hash to the dag.
func (i *importer) AddCommit(hash plumbing.Hash) (cid.Cid, error) {
	commit, err := i.repo.CommitObject(hash)
	if err != nil {
		return cid.Cid{}, err
	}

	var parents []cid.Cid
	for _, h := range commit.ParentHashes {
		parent, err := i.AddCommit(h)
		if err != nil {
			return cid.Cid{}, err
		}

		parents = append(parents, parent)
	}

	tree, err := i.AddTree(commit.TreeHash)
	if err != nil {
		return cid.Cid{}, err
	}

	mcommit := data.NewCommit(tree.Cid(), commit.Message, parents...)
	mcommit.Metadata["git_hash"] = hash.String()
	mcommit.Metadata["git_author_name"] = commit.Author.Name
	mcommit.Metadata["git_author_email"] = commit.Author.Email
	mcommit.Metadata["git_author_date"] = commit.Author.When.Format(DateFormat)
	mcommit.Metadata["git_committer_name"] = commit.Committer.Name
	mcommit.Metadata["git_committer_email"] = commit.Committer.Email
	mcommit.Metadata["git_committer_date"] = commit.Committer.When.Format(DateFormat)

	return data.AddCommit(i.ctx, i.dag, mcommit)
}

// AddTree adds the tree with the given hash to the dag.
func (i *importer) AddTree(hash plumbing.Hash) (ipld.Node, error) {
	tree, err := i.repo.TreeObject(hash)
	if err != nil {
		return nil, err
	}

	dir := ufsio.NewDirectory(i.dag)
	for _, entry := range tree.Entries {
		subnode, err := i.AddTreeEntry(entry)
		if err != nil {
			return nil, err
		}

		if err := dir.AddChild(i.ctx, entry.Name, subnode); err != nil {
			return nil, err
		}
	}

	node, err := dir.GetNode()
	if err != nil {
		return nil, err
	}

	if err := i.dag.Add(i.ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

// AddTreeEntry adds the tree entry with the given hash to the dag.
func (i *importer) AddTreeEntry(entry object.TreeEntry) (ipld.Node, error) {
	if entry.Mode == filemode.Dir {
		return i.AddTree(entry.Hash)
	}

	blob, err := i.repo.BlobObject(entry.Hash)
	if err != nil {
		return nil, err
	}

	r, err := blob.Reader()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if entry.Mode.IsFile() {
		return unixfs.Chunk(i.ctx, i.dag, r)
	}

	target, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	node := merkledag.NodeWithData(target)
	if err := i.dag.Add(i.ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}
