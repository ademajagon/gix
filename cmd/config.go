package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ademajagon/gix/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Gix configuration",
}

var keyProvider string

var setKeyCmd = &cobra.Command{
	Use:   "set-key",
	Short: "Set your API key for a provider",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter your %s API key: ", keyProvider)
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if key == "" {
			fmt.Println("API key cannot be empty.")
			return
		}

		// Load existing config to preserve other fields
		cfg, _ := config.Load()

		switch keyProvider {
		case "gemini":
			cfg.GeminiKey = key
		default:
			cfg.OpenAIKey = key
		}

		if err := config.Save(cfg); err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			return
		}

		fmt.Printf("%s API key saved successfully.\n", keyProvider)
	},
}

var setProviderCmd = &cobra.Command{
	Use:   "set-provider <openai|gemini>",
	Short: "Set the default AI provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToLower(strings.TrimSpace(args[0]))

		if name != "openai" && name != "gemini" {
			fmt.Fprintf(os.Stderr, "unknown provider %q (supported: openai, gemini)\n", name)
			os.Exit(1)
		}

		// Load existing config to preserve other fields
		cfg, _ := config.Load()
		cfg.Provider = name

		if err := config.Save(cfg); err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			return
		}

		fmt.Printf("Default provider set to %q.\n", name)
	},
}

func init() {
	setKeyCmd.Flags().StringVar(&keyProvider, "provider", "openai", "Provider to set the key for (openai, gemini)")
	configCmd.AddCommand(setKeyCmd)
	configCmd.AddCommand(setProviderCmd)
	rootCmd.AddCommand(configCmd)
}
