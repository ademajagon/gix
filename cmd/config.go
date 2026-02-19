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
	Short: "Manage gix configuration",
}

var setKeyCmd = &cobra.Command{
	Use:   "set-key",
	Short: "Set the API key for a provider",
	Long:  "Set the API key for an AI provider.",
	RunE:  runSetKey,
}

var setProviderCmd = &cobra.Command{
	Use:       "set-provider <openai|gemini>",
	Short:     "Set the default AI provider",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"openai", "gemini"},
	RunE:      runSetProvider,
}

var keyProvider string

func init() {
	setKeyCmd.Flags().StringVar(&keyProvider, "provider", "openai", "Provider to set the key for (openai, gemini)")

	configCmd.AddCommand(setKeyCmd)
	configCmd.AddCommand(setProviderCmd)
	rootCmd.AddCommand(configCmd)
}

func runSetKey(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter your %s API key: ", keyProvider)

	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	cfg, _ := config.Load()

	switch keyProvider {
	case "gemini":
		cfg.GeminiKey = key
	case "openai":
		cfg.OpenAIKey = key
	default:
		return fmt.Errorf("unknown provider", keyProvider)
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	path, _ := config.Path()
	fmt.Printf("%s API key saved to %s\n", keyProvider, path)
	return nil
}

func runSetProvider(_ *cobra.Command, args []string) error {
	name := strings.ToLower(strings.TrimSpace(args[0]))
	if name != "openai" && name != "gemini" {
		return fmt.Errorf("unknown provider", name)
	}

	cfg, _ := config.Load()
	cfg.Provider = name

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Printf("Default provider set to %q\n", name)
	return nil
}
