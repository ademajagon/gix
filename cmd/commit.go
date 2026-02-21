package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ademajagon/gix/config"
	"github.com/ademajagon/gix/internal/git"
	"github.com/ademajagon/gix/provider"
	"github.com/ademajagon/gix/utils"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate an AI commit message for staged changes",
	Long: `Generate a conventional commit message for your staged git diff using AI.
  [Enter]   accept and commit
  e         open in $EDITOR
  r         regenerate (ask the AI again)
  c         cancel`,
	RunE: runCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func runCommit(cmd *cobra.Command, _ []string) error {
	if !git.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	hasStaged, err := git.HasStagedChanges()
	if err != nil {
		return fmt.Errorf("checking staged changes: %w", err)
	}
	if !hasStaged {
		fmt.Fprintln(os.Stderr, "nothing to commit (no staged changes)")
		return nil
	}

	diff, err := git.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("reading diff: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	p, err := provider.New(cfg.ResolveProvider(), cfg.APIKey())
	if err != nil {
		return err
	}

	spinner := utils.NewSpinner()
	spinner.Start()
	suggestion, err := p.GenerateCommitMessage(diff)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("AI provider: %w", err)
	}

	finalMessage, err := promptMessage(suggestion, diff, p)
	if err != nil {
		fmt.Fprintln(os.Stderr, "aborted")
		return nil
	}

	gitCmd := exec.Command("git", "commit", "-m", finalMessage)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}

	return nil
}

// promptMessage runs the accept/edit/regenerate/cancel
func promptMessage(initial, diff string, p provider.AIProvider) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	msg := initial

	displayMessage(initial)

	for {
		fmt.Print("[Enter] commit  [e]dit  [r]egen  [c]ancel: ")

		raw, _ := reader.ReadString('\n')
		input := strings.TrimSpace(strings.ToLower(raw))

		switch input {
		case "":
			return msg, nil
		case "e":
			edited := utils.EditInEditor(msg)
			if edited == "" {
				fmt.Fprintln(os.Stderr, "commit message cannot be empty")
				continue
			}
			msg = edited
			displayMessage(msg)
		case "r":
			spinner := utils.NewSpinner()
			spinner.Start()
			newMsg, err := p.GenerateCommitMessage(diff)
			spinner.Stop()
			if err != nil {
				fmt.Fprintf(os.Stderr, "regen failed: %v\n", err)
				continue
			}
			msg = newMsg
			displayMessage(msg)
		case "c":
			return "", fmt.Errorf("cancelled")

		default:
			fmt.Fprintln(os.Stderr, "invalid input")
		}
	}
}

func displayMessage(msg string) {
	fmt.Print("\n> ")
	utils.TypingEffect(msg, 5*time.Millisecond)
	fmt.Println()
}
