package provider

import "time"

const (
	// Defaults, can be configured via gix config
	ollamaDefaultBaseURL   = "http://localhost:11434"
	ollamaDefaultChatModel = "llama3.2"
	// nomic-embed-text is a dedicated embedding model for Ollama.
	// Do NOT use llama3.2 for embeddings since it's a generative model
	// and will produce poor results for gix split's cosine clustering.
	ollamaDefaultEmbedModel = "nomic-embed-text"
)

type Ollama struct{ *chatClient }

func NewOllama(baseURL, chatModel, embedModel string) *Ollama {
	if baseURL == "" {
		baseURL = ollamaDefaultBaseURL
	}
	if chatModel == "" {
		chatModel = ollamaDefaultChatModel
	}
	if embedModel == "" {
		embedModel = ollamaDefaultEmbedModel
	}

	return &Ollama{newChatClient(
		baseURL+"/v1/chat/completions",
		baseURL+"/v1/embeddings",
		chatModel,
		embedModel,
		"", // ollama does not require api key
		60*time.Second,
	)}
}
