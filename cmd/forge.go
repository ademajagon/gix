package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ademajagon/gix/config"
	"github.com/ademajagon/gix/git"
	"github.com/ademajagon/gix/openai"
	"github.com/ademajagon/gix/utils"
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

		cfg, _ := config.Load()
		diff, _ := git.GetFullDiff()
		hunks := git.ExtractHunks(diff)

		var contents []string
		for _, h := range hunks {
			contents = append(contents, h.Content)
		}

		vectors, err := openai.EmbedHunks(cfg.OpenAIKey, contents)
		if err != nil {
			fmt.Fprintf(os.Stderr, "embedding error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Got %d vectors\n", len(vectors))

		clusters := utils.ClusterGroups(vectors, 0.80)

		for i, group := range clusters {
			var groupHunks []string

			for _, idx := range group {
				groupHunks = append(groupHunks, hunks[idx].Content)
			}

			msg, err := openai.GenerateCommitMessageForGroup(cfg.OpenAIKey, groupHunks)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to generate commit message for group %d: %v\n", i+1, err)
				continue
			}

			fmt.Printf("\nGroup %d:\n", i+1)
			utils.TypingEffect("Commit message: "+msg, 5*time.Millisecond)

			patch := ""
			for _, h := range groupHunks {
				patch += h + "\n"
			}

			if err := utils.ApplyPatchToIndex(patch); err != nil {
				utils.TypingEffect("Patch failed. Skipping group.", 5*time.Millisecond)
				continue
			}

			if err := utils.CommitWithMessage(msg); err != nil {
				utils.TypingEffect("Commit failed. Reseting index.", 5*time.Millisecond)
				_ = utils.ResetIndex()
				continue
			}

			utils.TypingEffect("Commited successfully!", 5*time.Millisecond)
			_ = utils.ResetIndex()
		}
	},
}

func init() {
	rootCmd.AddCommand(forgeCmd)
}
