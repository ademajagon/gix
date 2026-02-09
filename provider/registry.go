package provider

import "fmt"

const (
	ProviderOpenAI = "openai"
	ProviderGemini = "gemini"
)

// New creates an AIProvider for the given provider name.
func New(name string, apiKey string) (AIProvider, error) {
	switch name {
	case ProviderOpenAI:
		return NewOpenAI(apiKey), nil
	case ProviderGemini:
		return NewGemini(apiKey), nil
	default:
		return nil, fmt.Errorf("unknown provider %q (supported: openai, gemini)", name)
	}
}
