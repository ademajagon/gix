package cmd

import (
	"fmt"
	"os"

	"github.com/ademajagon/gix/git"
	"github.com/spf13/cobra"
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Suggest a split of the current staged diff into multiple semantic commits",
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintf(os.Stderr, "fatal: not a git repository")
			os.Exit(1)
		}

		hasStaged, err := git.HasStagedChanges()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to check staged changes: %v\n", err)
			os.Exit(1)
		}

		if !hasStaged {
			fmt.Fprintln(os.Stderr, "nothing to split (no staged changes)")
			os.Exit(0)
		}

		hunks, err := git.ParseHunks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to parse hunks: %v\n", err)
			os.Exit(1)
		}

		if len(hunks) == 0 {
			fmt.Println("No hunks found in staged diff.")
			return
		}

		fmt.Printf("\nFound %d hunks:\n\n", len(hunks))
		for i, h := range hunks {
			fmt.Printf("Hunk %d:\n", i+1)
			fmt.Printf("File:   %s\n", h.FilePath)
			fmt.Printf("Header: %s\n", h.Header)
			fmt.Println("Body:")
			fmt.Println(h.Body)
			fmt.Println("--------------------------------------------------")
		}
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)
}
