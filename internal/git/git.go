// Package git contains methods for importing Git repositories.
package git

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	ufs "github.com/ipfs/go-unixfs"
	ufsio "github.com/ipfs/go-unixfs/io"

	"github.com/multiverse-vcs/go-multiverse/pkg/dag"
	mobject "github.com/multiverse-vcs/go-multiverse/pkg/object"
)

// DateFormat is the format to store dates in.
const DateFormat = "Mon Jan 02 15:04:05 2006 -0700"

// importer adds objects from a git repo to a dag.
type importer struct {
	ctx      context.Context
	dag      ipld.DAGService
	name     string
	repo     *git.Repository
	objects  map[string]cid.Cid
	branches map[string]cid.Cid
	tags     map[string]cid.Cid
}

// ImportFromURL is a helper to import a git repo from a url.
func ImportFromURL(ctx context.Context, dag ipld.DAGService, name, url string) (cid.Cid, error) {
	dir := filepath.Join(os.TempDir(), "multi_git_import_"+name)
	defer os.RemoveAll(dir)

	opts := git.CloneOptions{
		URL: url,
	}

	repo, err := git.PlainClone(dir, true, &opts)
	if err != nil {
		return cid.Cid{}, err
	}

	return NewImporter(ctx, dag, repo, name).AddRepository()
}

// ImportFromFS is a helper to import a git repo from a directory.
func ImportFromFS(ctx context.Context, dag ipld.DAGService, name, dir string) (cid.Cid, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return cid.Cid{}, err
	}

	return NewImporter(ctx, dag, repo, name).AddRepository()
}

// NewImporter returns an importer for the given repo.
func NewImporter(ctx context.Context, dag ipld.DAGService, repo *git.Repository, name string) *importer {
	return &importer{
		ctx:      ctx,
		dag:      dag,
		name:     name,
		repo:     repo,
		objects:  make(map[string]cid.Cid),
		branches: make(map[string]cid.Cid),
		tags:     make(map[string]cid.Cid),
	}
}

// AddRepository adds all branches and tags to the dag.
func (i *importer) AddRepository() (cid.Cid, error) {
	head, err := i.repo.Head()
	if err != nil {
		return cid.Cid{}, err
	}

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

	defaultBranch := string(head.Name())
	defaultBranch = path.Base(defaultBranch)

	mrepo := mobject.NewRepository()
	mrepo.Branches = i.branches
	mrepo.Tags = i.tags
	mrepo.DefaultBranch = defaultBranch

	return mobject.AddRepository(i.ctx, i.dag, mrepo)
}

// AddBranch adds the branch with the given ref to the dag.
func (i *importer) AddBranch(ref *plumbing.Reference) error {
	id, err := i.AddCommit(ref.Hash())
	if err != nil {
		return err
	}

	name := string(ref.Name())
	name = path.Base(name)

	i.branches[name] = id
	return nil
}

// AddTag adds the tag with the given ref to the dag.
func (i *importer) AddTag(ref *plumbing.Reference) error {
	id, ok := i.objects[ref.Hash().String()]
	if !ok {
		return nil
	}

	name := string(ref.Name())
	name = path.Base(name)

	i.tags[name] = id
	return nil
}

// AddCommit adds the commit with the given hash to the dag.
func (i *importer) AddCommit(hash plumbing.Hash) (cid.Cid, error) {
	if id, ok := i.objects[hash.String()]; ok {
		return id, nil
	}

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

	mcommit := mobject.NewCommit()
	mcommit.Tree = tree.Cid()
	mcommit.Message = commit.Message
	mcommit.Parents = parents
	mcommit.Date = commit.Committer.When
	mcommit.Metadata["git_hash"] = hash.String()
	mcommit.Metadata["git_author_name"] = commit.Author.Name
	mcommit.Metadata["git_author_email"] = commit.Author.Email
	mcommit.Metadata["git_committer_name"] = commit.Committer.Name
	mcommit.Metadata["git_committer_email"] = commit.Committer.Email

	id, err := mobject.AddCommit(i.ctx, i.dag, mcommit)
	if err != nil {
		return cid.Cid{}, err
	}

	i.objects[hash.String()] = id
	return id, nil
}

// AddTree adds the tree with the given hash to the dag.
func (i *importer) AddTree(hash plumbing.Hash) (ipld.Node, error) {
	if id, ok := i.objects[hash.String()]; ok {
		return i.dag.Get(i.ctx, id)
	}

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

	i.objects[hash.String()] = node.Cid()
	return node, nil
}

// AddTreeEntry adds the tree entry with the given hash to the dag.
func (i *importer) AddTreeEntry(entry object.TreeEntry) (ipld.Node, error) {
	switch entry.Mode {
	case filemode.Dir:
		return i.AddTree(entry.Hash)
	case filemode.Submodule:
		return ufs.EmptyDirNode(), nil
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
		return dag.Chunk(i.ctx, i.dag, r)
	}

	target, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	node := merkledag.NodeWithData(target)
	if err := i.dag.Add(i.ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}
