package repo

import (
	"context"

	"github.com/multiverse-vcs/go-multiverse/pkg/object"
	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
)

// ListArgs contains the args.
type ListArgs struct{}

// ListReply contains the reply
type ListReply struct {
	// Repositories is a map of repositories.
	Repositories map[remote.Path]*object.Repository `json:"repositories"`
}

// List returns a list of repositories.
func (s *Service) List(args *ListArgs, reply *ListReply) error {
	ctx := context.Background()

	author := s.Config.Author
	peerID := s.Peer.Host.ID()

	repos := make(map[remote.Path]*object.Repository)
	for name, id := range author.Repositories {
		repo, err := object.GetRepository(ctx, s.Peer.DAG, id)
		if err != nil {
			return err
		}

		path := remote.NewPath(peerID, name)
		repos[path] = repo
	}

	reply.Repositories = repos
	return nil
}
