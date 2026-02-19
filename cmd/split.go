package cmd

import (
	"fmt"
	"os"

	"github.com/ademajagon/gix/config"
	"github.com/ademajagon/gix/internal/git"
	"github.com/ademajagon/gix/provider"
	"github.com/ademajagon/gix/split"
	"github.com/ademajagon/gix/utils"
	"github.com/spf13/cobra"
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "[BETA] Split staged changes into semantic atomic commits",
	Long:  "[BETA] Split staged changes into multiple semantic commits using AI.",
	RunE:  runSplit,
}

func init() {
	rootCmd.AddCommand(splitCmd)
}

func runSplit(_ *cobra.Command, _ []string) error {
	if !git.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	hasStaged, err := git.HasStagedChanges()
	if err != nil {
		return fmt.Errorf("checking staged changes: %w", err)
	}
	if !hasStaged {
		fmt.Fprintln(os.Stderr, "nothing to split (no staged changes)")
		return nil
	}

	hunks, err := git.ParseHunks()
	if err != nil {
		return fmt.Errorf("parsing hunks: %w", err)
	}
	if len(hunks) == 0 {
		fmt.Fprintln(os.Stderr, "no hunks found in staged diff")
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	p, err := provider.New(cfg.ResolveProvider(), cfg.APIKey())
	if err != nil {
		return err
	}

	fmt.Printf("[BETA] Analysing %d hunk(s)…\n", len(hunks))

	spinner := utils.NewSpinner()
	spinner.Start()
	groups, err := split.ClusterHunks(p, hunks)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("clustering hunks: %w", err)
	}

	if len(groups) == 0 {
		fmt.Fprintln(os.Stderr, "no semantic groups found")
		return nil
	}

	fmt.Printf("\nProposed %d commit(s):\n", len(groups))
	for i, g := range groups {
		fmt.Printf("  %d. %s (%d hunk(s))\n", i+1, g.Message, len(g.Hunks))
	}

	fmt.Print("\nApply these commits? [y/N]: ")
	var answer string
	fmt.Scanln(&answer)
	if answer != "y" && answer != "Y" {
		fmt.Fprintln(os.Stderr, "aborted")
		return nil
	}

	if err := split.ApplyGroups(groups); err != nil {
		return fmt.Errorf("applying commits: %w", err)
	}

	fmt.Printf("\nCreated %d commit(s).\n", len(groups))
	return nil
}
