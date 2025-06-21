package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message using AI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating commit message... (placeholder)")
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
