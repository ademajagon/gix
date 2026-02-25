package provider

import (
	"fmt"

	"github.com/ademajagon/gix/config"
)

const (
	ProviderOpenAI = "openai"
	ProviderGemini = "gemini"
	ProviderOllama = "ollama"
)

// New returns an AIProvider for the given name and API key.
func New(name, apiKey string) (AIProvider, error) {
	switch name {
	case ProviderOpenAI:
		if apiKey == "" {
			return nil, fmt.Errorf("API key is required for provider %q, run `gix config set-key`", name)
		}
		return NewOpenAI(apiKey), nil
	case ProviderGemini:
		if apiKey == "" {
			return nil, fmt.Errorf("API key is required for provider %q, run `gix config set-key`", name)
		}
		return NewGemini(apiKey), nil
	case ProviderOllama:
		return NewOllama("", "", ""), nil
	default:
		return nil, fmt.Errorf("unknown provider %q (supported: openai, gemini, ollama)", name)
	}
}

func NewFromConfig(cfg config.Config) (AIProvider, error) {
	switch cfg.ResolveProvider() {
	case ProviderOllama:
		return NewOllama(cfg.OllamaBaseURL, cfg.OllamaChatModel, cfg.OllamaEmbedModel), nil
	default:
		return New(cfg.ResolveProvider(), cfg.APIKey())
	}
}
