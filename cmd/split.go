package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ademajagon/gix/config"
	"github.com/ademajagon/gix/git"
	"github.com/ademajagon/gix/semantics"
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
			fmt.Fprintf(os.Stderr, "error: checking staged changes: %v\n", err)
			os.Exit(1)
		}

		if !hasStaged {
			fmt.Fprintln(os.Stderr, "nothing to split (no staged changes)")
			os.Exit(0)
		}

		hunks, err := git.ParseHunks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing hunks: %v\n", err)
			os.Exit(1)
		}

		if len(hunks) == 0 {
			fmt.Println("No hunks found in staged diff.")
			return
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
			os.Exit(1)
		}

		groups, err := semantics.ClusterHunks(cfg.OpenAIKey, hunks)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error clustering hunks: %v\n", err)
			os.Exit(1)
		}

		for i, group := range groups {
			fmt.Printf("Commit %d: %s\n", i+1, group.Message)
			for _, h := range group.Hunks {
				fmt.Printf("- %s: %s\n", h.FilePath, strings.TrimSpace(h.Header))
			}
			fmt.Println("--------------------------------------------------")
		}
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)
}
