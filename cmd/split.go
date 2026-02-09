package cmd

import (
	"fmt"
	"os"

	"github.com/ademajagon/gix/config"
	"github.com/ademajagon/gix/git"
	"github.com/ademajagon/gix/provider"
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

		p, err := provider.New(cfg.ResolveProvider(), cfg.APIKey())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		groups, err := semantics.ClusterHunks(p, hunks)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error clustering hunks: %v\n", err)
			os.Exit(1)
		}

		if len(groups) == 0 {
			fmt.Fprintln(os.Stderr, "nothing to split (no semantic groups found)")
			os.Exit(0)
		}

		if err := semantics.ApplyGroups(groups); err != nil {
			fmt.Fprintf(os.Stderr, "error applying commits: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)
}
