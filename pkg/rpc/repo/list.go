package repo

import (
	"sort"
)

// ListArgs contains the args.
type ListArgs struct{}

// ListReply contains the reply
type ListReply struct {
	// Repositories is a list of repositories.
	Repositories []string `json:"repositories"`
}

// List returns a list of repositories.
func (s *Service) List(args *ListArgs, reply *ListReply) error {
	var names []string
	for name := range s.Config.Author.Repositories {
		names = append(names, name)
	}
	sort.Strings(names)

	reply.Repositories = names
	return nil
}
