package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

		spinner := utils.NewSpinner()
		spinner.Start()

		suggestion, err := openai.GenerateCommitMessage(cfg.OpenAIKey, diff)

		spinner.Stop()

		if err != nil {
			fmt.Fprintf(os.Stderr, "OpenAI error: %v\n", err)
			os.Exit(1)
		}

		firstLine := suggestion
		if i := indexOfNewline(suggestion); i != -1 {
			firstLine = suggestion[:i]
		}

		fullCmd := fmt.Sprintf("git commit -m %q", firstLine)
		utils.TypingEffect(fullCmd, 5*time.Millisecond)

		fmt.Print("\n\nPress [Enter] to commit, or Ctrl+C to cancel...")
		_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')

		genCmd := exec.Command("git", "commit", "-m", suggestion)
		genCmd.Stdout = os.Stdout
		genCmd.Stderr = os.Stderr
		if err := genCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Commit failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func indexOfNewline(s string) int {
	for i, c := range s {
		if c == '\n' || c == '\r' {
			return i
		}
	}
	return -1
}
