package cmd

import (
	"context"
	"os"

	http "github.com/ipfs/go-ipfs-http-client"

	"github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
	"github.com/spf13/cobra"
)

var message string

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes.",
	Long: `Commit changes.`,
	RunE: executeCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Description of changes made")
}

func executeCommit(cmd *cobra.Command, args []string) error {
	api, err := http.NewLocalApi()
	if err != nil {
		return err
	}

	dir, err := os.Open(".")
	if err != nil {
		return err
	}

	stat, err := dir.Stat()
	if err != nil {
		return err
	}

	node, err := files.NewSerialFile(".", true, stat)
	if err != nil {
		return err
	}

	path, err := api.Unixfs().Add(context.TODO(), node)
	if err != nil {
		return err
	}

	commit := make(map[string]interface{})
	commit["message"] = message
	commit["changes"] = path.Cid()

	dag, err := cbornode.WrapObject(commit, multihash.SHA2_256, -1)
	if err != nil {
		return err
	}

	return api.Dag().Pinning().Add(context.TODO(), dag)
}