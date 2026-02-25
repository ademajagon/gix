package provider

import "time"

const (
	openaiChatURL    = "https://api.openai.com/v1/chat/completions"
	openaiEmbedURL   = "https://api.openai.com/v1/embeddings"
	openaiChatModel  = "gpt-4o"
	openaiEmbedModel = "text-embedding-3-small"
)

type OpenAI struct{ *chatClient }

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{newChatClient(
		openaiChatURL,
		openaiEmbedURL,
		openaiChatModel,
		openaiEmbedModel,
		apiKey,
		20*time.Second,
	)}
}
