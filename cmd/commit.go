package cmd

import (
	"fmt"
	"os"

	"github.com/ademajagon/toka/git"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message using AI",
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "Not inside a Git repository.")
			os.Exit(1)
		}

		hasStaged, err := git.HasStagedChanges()
		if err != nil || !hasStaged {
			fmt.Fprintln(os.Stderr, "No staged changes to commit.")
			os.Exit(1)
		}

		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get Git diff: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Staged changes found.")
		fmt.Println("--- DIFF BEGIN ---")
		fmt.Println(diff)
		fmt.Println("--- DIFF END ---")
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
