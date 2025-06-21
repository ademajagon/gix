package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "toka",
	Short: "Toka is an AI-powered commit assistant.",
	Long:  `Toka is a CLI tool that suggests Git commit messages using AI based on your staged changes.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
