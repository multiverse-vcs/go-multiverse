package repo

import (
	cid "github.com/ipfs/go-cid"
)

// ListArgs contains the args.
type ListArgs struct{}

// ListReply contains the reply
type ListReply struct {
	// Repositories is a map of repositories.
	Repositories map[string]cid.Cid `json:"repositories"`
}

// List returns a list of repositories.
func (s *Service) List(args *ListArgs, reply *ListReply) error {
	reply.Repositories = s.Config.Author.Repositories
	return nil
}
