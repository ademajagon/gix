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
	Long: `Set the API key for an AI provider.

Examples:
  gix config set-key                      # set OpenAI key (default)
  gix config set-key --provider gemini    # set Gemini key`,
	RunE: runSetKey,
}

var setProviderCmd = &cobra.Command{
	Use:       "set-provider <openai|gemini>",
	Short:     "Set the default AI provider",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"openai", "gemini", "ollama"},
	RunE:      runSetProvider,
}

var setUpdateCheckCmd = &cobra.Command{
	Use:   "update-check <on|off>",
	Short: "Enable or disable background version update checks",
	Long: `Enable or disable background version update checks.

gix checks for newer versions in the background after each command and prints
a notice if one is available. The check never blocks the primary command and
results are cached for 48 hours.

To disable permanently:
  gix config update-check off

You can also set GIX_CHECKPOINT_DISABLE=1 in your environment to disable
for a single session without changing the config file.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"on", "off"},
	RunE:      runSetUpdateCheck,
}

var setOllamaURLCmd = &cobra.Command{
	Use:   "set-ollama-url <url>",
	Short: "Set the Ollama base URL (default: http://localhost:11434)",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetOllamaURL,
}

var setOllamaModelCmd = &cobra.Command{
	Use:   "set-ollama-model <chat-model> <embed-model>",
	Short: "Set the Ollama chat and embed models",
	Long: `Set the Ollama models used for commit messages and gix split.

Example:
  gix config set-ollama-model llama3.2 nomic-embed-text`,
	Args: cobra.ExactArgs(2),
	RunE: runSetOllamaModel,
}

var keyProvider string

func init() {
	setKeyCmd.Flags().StringVar(&keyProvider, "provider", "openai", "Provider to set the key for (openai, gemini)")

	configCmd.AddCommand(setKeyCmd)
	configCmd.AddCommand(setProviderCmd)
	configCmd.AddCommand(setUpdateCheckCmd)
	configCmd.AddCommand(setOllamaURLCmd)
	configCmd.AddCommand(setOllamaModelCmd)
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

	cfg, _ := config.Load() // ignore error, file may not exist yet

	switch keyProvider {
	case "gemini":
		cfg.GeminiKey = key
	case "openai":
		cfg.OpenAIKey = key
	default:
		return fmt.Errorf("unknown provider %q", keyProvider)
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
	if name != "openai" && name != "gemini" && name != "ollama" {
		return fmt.Errorf("unknown provider %q", name)
	}

	cfg, _ := config.Load()
	cfg.Provider = name

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Printf("Default provider set to %q\n", name)
	return nil
}

func runSetUpdateCheck(_ *cobra.Command, args []string) error {
	val := strings.ToLower(strings.TrimSpace(args[0]))
	if val != "on" && val != "off" {
		return fmt.Errorf("expected 'on' or 'off', got %q", val)
	}

	cfg, _ := config.Load()
	cfg.DisableUpdateCheck = val == "off"

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	if cfg.DisableUpdateCheck {
		fmt.Println("Update checks disabled.")
		fmt.Println("You can also set GIX_CHECKPOINT_DISABLE=1 for a one-off session disable.")
	} else {
		fmt.Println("Update checks enabled.")
	}

	return nil
}

func runSetOllamaURL(_ *cobra.Command, args []string) error {
	cfg, _ := config.Load()
	cfg.OllamaBaseURL = strings.TrimSpace(args[0])
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("Ollama base URL set to %q\n", cfg.OllamaBaseURL)
	return nil
}

func runSetOllamaModel(_ *cobra.Command, args []string) error {
	cfg, _ := config.Load()
	cfg.OllamaChatModel = strings.TrimSpace(args[0])
	cfg.OllamaEmbedModel = strings.TrimSpace(args[1])
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Printf("Ollama models set — chat: %q, embed: %q\n", cfg.OllamaChatModel, cfg.OllamaEmbedModel)
	return nil
}
