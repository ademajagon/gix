package cmd

import (
	"fmt"
	"os"

	"github.com/ademajagon/gix/config"
	"github.com/ademajagon/gix/git"
	"github.com/ademajagon/gix/openai"
	"github.com/spf13/cobra"
)

var forgeCmd = &cobra.Command{
	Use:   "forge",
	Short: "Forge structured commits from chaotic changes using AI",
	Long:  `Gix Forge analyzes your staged or working changes and uses AI to help split them into well-structured commits - grouped by tickets, topic or logic.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "fatal: not a git repository")
			os.Exit(1)
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: config not loaded")
			os.Exit(1)
		}

		diff, err := git.GetFullDiff()
		if err != nil {
			fmt.Fprintln(os.Stderr, "nothing to forge (no unstaged changes)")
			os.Exit(0)
		}

		fmt.Println("Analyzing your changes...")

		grouped, err := openai.GroupDiffIntoCommits(cfg.OpenAIKey, diff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: OpenAI failed: %v\n", err)
			os.Exit(1)
		}

		if len(grouped) == 0 {
			fmt.Println("No groups were returned.")
			return
		}

		fmt.Println("\n Suggested commit groupings:")

		for i, g := range grouped {
			fmt.Printf("  %d. %s\n", i+1, g)
		}
	},
}

func init() {
	rootCmd.AddCommand(forgeCmd)
}
