package cmd

import (
	"context"
	"path/filepath"

	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/repo"
)

var cloneCmd = &cobra.Command{
	Use:          "clone [remote] [local]",
	Short:        "Copy an existing repository.",
	Long:         `Copy an existing repository.`,
	Args:         cobra.ExactArgs(2),
	SilenceUsage: true,
	RunE:         executeClone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func executeClone(cmd *cobra.Command, args []string) error {
	api, err := httpapi.NewLocalApi()
	if err != nil {
		return err
	}

	remote, err := api.ResolvePath(context.TODO(), path.Join(path.New(args[0]), "tree"))
	if err != nil {
		return err
	}

	f, err := api.Unixfs().Get(context.TODO(), remote)
	if err != nil {
		return err
	}

	local, err := filepath.Abs(args[1])
	if err != nil {
		return err
	}

	if err := files.WriteTo(f, local); err != nil {
		return err
	}

	r := repo.Repo{Path: local, Head: remote.Root()}
	return r.Write()
}
