package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yondero/multiverse/repo"
)

var initCmd = &cobra.Command{
	Use: "init [path]",
	Short: "Initialize a Multiverse repository.",
	Long: `Initialize a Multiverse repository.`,
	Args: cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: executeInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func executeInit(cmd *cobra.Command, args []string) error {
	var path = "."
	if len(args) > 0 {
		path = args[0]
	}

	r, err := repo.NewRepo()
	if err != nil {
		return err
	}

	if dir, err := r.Dir(path); err == nil {
		return fmt.Errorf("Repo exists at %s", dir)
	}

	if err := r.Init(path); err != nil {
		return err
	}

	dir, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	fmt.Println("Repo initalized at", dir)
	return nil
}