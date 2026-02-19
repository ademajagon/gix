package provider

import (
	"fmt"
)

const (
	POpenAI = "openai"
	PGemini = "gemini"
)

// New returns an AIProvider for the given name and API key.
func New(name, apiKey string) (AIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for provider %q — run `gix config set-key`", name)
	}

	switch name {
	case POpenAI:
		return NewOpenAI(apiKey), nil
	case PGemini:
		return NewGemini(apiKey), nil
	default:
		return nil, fmt.Errorf("unknown provider %q (supported: openai, gemini)", name)
	}
}
