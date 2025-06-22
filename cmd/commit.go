package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ademajagon/toka/config"
	"github.com/ademajagon/toka/git"
	"github.com/ademajagon/toka/openai"
	"github.com/ademajagon/toka/utils"
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

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to load configuration:", err)
			os.Exit(1)
		}

		suggestion, err := openai.GenerateCommitMessage(cfg.OpenAIKey, diff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "OpenAI error: %v\n", err)
			os.Exit(1)
		}

		utils.TypingEffect(suggestion, 5*time.Millisecond)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
