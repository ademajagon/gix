package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ademajagon/gix/config"
	"github.com/ademajagon/gix/git"
	"github.com/ademajagon/gix/openai"
	"github.com/ademajagon/gix/utils"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate a commit message using AI",
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitRepo() {
			fmt.Fprintln(os.Stderr, "fatal: not a git repository")
			os.Exit(1)
		}

		hasStaged, err := git.HasStagedChanges()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to check staged changes: %v\n", err)
			os.Exit(1)
		}
		if !hasStaged {
			fmt.Fprintln(os.Stderr, "nothing to commit (no staged changes)")
			os.Exit(0)
		}

		diff, err := git.GetStagedDiff()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to read diff: %v\n", err)
			os.Exit(1)
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: config not loaded: %v\n", err)
			fmt.Fprintln(os.Stderr, "hint: run `gix config set-key` to set your OpenAI key")
			os.Exit(1)
		}

		spinner := utils.NewSpinner()
		spinner.Start()
		suggestion, err := openai.GenerateCommitMessage(cfg.OpenAIKey, diff)
		spinner.Stop()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: OpenAI failed: %v\n", err)
			os.Exit(1)
		}

		finalMessage, err := promptCommitMessage(suggestion, diff, cfg.OpenAIKey)
		if err != nil {
			os.Exit(0)
		}

		cmdExec := exec.Command("git", "commit", "-m", finalMessage)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: commit failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func promptCommitMessage(initial, diff, apiKey string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	msg := initial

	printWithEffect(msg)

	for {
		fmt.Print("[Enter] to commit  [e]dit  [r]egen  [c]ancel: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "":
			return msg, nil
		case "e":
			msg = utils.EditInEditor(msg)
			printWithEffect(msg)
		case "r":
			spinner := utils.NewSpinner()
			spinner.Start()
			newMsg, err := openai.GenerateCommitMessage(apiKey, diff)
			spinner.Stop()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: OpenAI failed: %v\n", err)
				continue
			}
			msg = newMsg
			printWithEffect(msg)
		case "c":
			fmt.Println("canceled")
			return "", fmt.Errorf("canceled")
		default:
			fmt.Println("invalid input")
		}
	}
}

func printWithEffect(msg string) {
	fmt.Print("\n> ")
	utils.TypingEffect(msg, 5*time.Millisecond)
	fmt.Println()
}
