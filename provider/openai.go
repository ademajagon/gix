package provider

import "time"

const (
	openaiChatURL    = "https://api.openai.com/v1/chat/completions"
	openaiEmbedURL   = "https://api.openai.com/v1/embeddings"
	defaultOpenAIChatModel  = "gpt-4o"
	defaultOpenAIEmbedModel = "text-embedding-3-small"
)

type OpenAI struct{ *chatClient }

func NewOpenAI(apiKey, chatModel, embedModel string) *OpenAI {
	if chatModel == "" {
		chatModel = defaultOpenAIChatModel
	}
	if embedModel == "" {
		embedModel = defaultOpenAIEmbedModel
	}
	return &OpenAI{newChatClient(
		openaiChatURL,
		openaiEmbedURL,
		chatModel,
		embedModel,
		apiKey,
		20*time.Second,
	)}
}
