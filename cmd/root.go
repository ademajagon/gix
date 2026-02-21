package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

// Version returns the current build version
func Version() string {
	return version
}

var showUpdateNotice func()

func SetUpdateNotice(fn func()) {
	showUpdateNotice = fn
}

var rootCmd = &cobra.Command{
	Use:          "gix",
	Short:        "AI powered git commit assistant",
	Long:         "gix helps you write clean conventional commit messages and split staged diffs into smaller commits.",
	Version:      version,
	SilenceUsage: true,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if showUpdateNotice != nil {
			showUpdateNotice()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
