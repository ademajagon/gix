package provider

import (
	"fmt"
)

const (
	ProviderOpenAI = "openai"
	ProviderGemini = "gemini"
)

// New returns an AIProvider for the given name and API key.
func New(name, apiKey string) (AIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for provider %q — run `gix config set-key`", name)
	}

	switch name {
	case ProviderOpenAI:
		return NewOpenAI(apiKey), nil
	case ProviderGemini:
		return NewGemini(apiKey), nil
	default:
		return nil, fmt.Errorf("unknown provider %q (supported: openai, gemini)", name)
	}
}
