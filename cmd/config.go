package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ademajagon/toka/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Toka configuration",
}

var setKeyCmd = &cobra.Command{
	Use:   "set-key <key>",
	Short: "Set your OpenAI API key",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your OpenAI API key: ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if key == "" {
			fmt.Println("API key cannot be empty.")
			return
		}

		cfg := config.Config{
			OpenAIKey: key,
		}
		if err := config.Save(cfg); err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			return
		}

		fmt.Println("OpenAI API key saved successfully.")
	},
}

func init() {
	configCmd.AddCommand(setKeyCmd)
	rootCmd.AddCommand(configCmd)
}
