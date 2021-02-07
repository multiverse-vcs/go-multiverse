package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// TagArgs contains the args.
type TagArgs struct {
	// Name is the name of the repo.
	Name string
	// Tag is the name of the tag.
	Tag string
	// Head is the CID of the repo head.
	Head cid.Cid
}

// TagReply contains the reply.
type TagReply struct {
	// Tags is a map of commit CIDs.
	Tags map[string]cid.Cid
}

// ListTags returns the repo tags.
func (s *Service) ListTags(args *TagArgs, reply *TagReply) error {
	ctx := context.Background()
	cfg := s.client.Config()

	id, ok := cfg.Author.Repositories[args.Name]
	if !ok {
		return errors.New("repository does not exist")
	}

	repo, err := data.GetRepository(ctx, s.client, id)
	if err != nil {
		return err
	}

	reply.Tags = repo.Tags
	return nil
}

// CreateTag creates a new tag.
func (s *Service) CreateTag(args *TagArgs, reply *TagReply) error {
	ctx := context.Background()
	cfg := s.client.Config()

	id, ok := cfg.Author.Repositories[args.Name]
	if !ok {
		return errors.New("repository does not exist")
	}

	repo, err := data.GetRepository(ctx, s.client, id)
	if err != nil {
		return err
	}

	if args.Tag == "" {
		return errors.New("name cannot be empty")
	}

	if _, ok := repo.Tags[args.Tag]; ok {
		return errors.New("tag already exists")
	}

	repo.Tags[args.Tag] = args.Head
	reply.Tags = repo.Tags

	id, err = data.AddRepository(ctx, s.client, repo)
	if err != nil {
		return err
	}

	cfg.Sequence++
	cfg.Author.Repositories[args.Name] = id

	if err := cfg.Save(); err != nil {
		return err
	}

	return s.client.Authors().Publish(ctx)
}

// DeleteTag deletes an existing tag.
func (s *Service) DeleteTag(args *TagArgs, reply *TagReply) error {
	ctx := context.Background()
	cfg := s.client.Config()

	id, ok := cfg.Author.Repositories[args.Name]
	if !ok {
		return errors.New("repository does not exist")
	}

	repo, err := data.GetRepository(ctx, s.client, id)
	if err != nil {
		return err
	}

	if args.Tag == "" {
		return errors.New("name cannot be empty")
	}

	if _, ok := repo.Tags[args.Tag]; !ok {
		return errors.New("tag does not exists")
	}

	delete(repo.Tags, args.Tag)
	reply.Tags = repo.Tags

	id, err = data.AddRepository(ctx, s.client, repo)
	if err != nil {
		return err
	}

	cfg.Sequence++
	cfg.Author.Repositories[args.Name] = id

	if err := cfg.Save(); err != nil {
		return err
	}

	return s.client.Authors().Publish(ctx)
}
