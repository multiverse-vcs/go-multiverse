package cmd

import (
	"os"

	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiverse-vcs/go-multiverse/core"
	"github.com/multiverse-vcs/go-multiverse/repo"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:          "switch [name]",
	Short:        "Change to a different branch.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE:         executeSwitch,
}

func init() {
	rootCmd.AddCommand(switchCmd)
}

func executeSwitch(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	name := args[0]

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err := repo.Open(cwd)
	if err != nil {
		return err
	}

	c, err := core.NewCore(ctx)
	if err != nil {
		return err
	}

	id, err := r.Branches.Head(name)
	if err != nil {
		return err
	}

	_, err = c.Reference(ctx, path.IpfsPath(id))
	if err != nil {
		return err
	}

	// TODO do not update tree if base is the same
	// if err := util.WriteTo(tree, r.Path); err != nil {
	// 	return err
	// }

	r.Base = id
	r.Branch = name
	return r.Write()
}
