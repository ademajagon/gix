package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	openaiChatURL    = "https://api.openai.com/v1/chat/completions"
	openaiChatModel  = "gpt-4o"
	openaiEmbedURL   = "https://api.openai.com/v1/embeddings"
	openaiEmbedModel = "text-embedding-3-small"
)

type OpenAI struct {
	apiKey     string
	httpClient *http.Client
}

func NewOpenAI(apiKey string) *OpenAI {
	return NewOpenAIWithClient(apiKey, &http.Client{Timeout: 20 * time.Second})
}

func NewOpenAIWithClient(apiKey string, client *http.Client) *OpenAI {
	return &OpenAI{apiKey: apiKey, httpClient: client}
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	Temperature float32         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

type openaiChoice struct {
	Message openaiMessage `json:"message"`
}

type openaiChatResponse struct {
	Choices []openaiChoice `json:"choices"`
	Error   *openaiError   `json:"error,omitempty"`
}

type openaiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Embedding types
type openaiEmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type openaiEmbedResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func (o *OpenAI) GenerateCommitMessage(diff string) (string, error) {
	payload := openaiChatRequest{
		Model: openaiChatModel,
		Messages: []openaiMessage{
			{Role: "system", Content: CommitMessageSystem},
			{Role: "user", Content: CommitMessageUser + diff},
		},
		Temperature: 0,
		MaxTokens:   128,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshalling request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, openaiChatURL, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	res, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("OpenAI request: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		var apiErr openaiChatResponse
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != nil {
			return "", fmt.Errorf("OpenAI API error: %s", apiErr.Error.Message)
		}
		return "", fmt.Errorf("OpenAI API error %s", res.Status)
	}

	var response openaiChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("decoding OpenAI response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", errors.New("OpenAI returned no choices")
	}

	return strings.TrimSpace(response.Choices[0].Message.Content), nil
}

func (o *OpenAI) GetEmbeddings(texts []string) ([][]float32, error) {
	payload := openaiEmbedRequest{
		Model: openaiEmbedModel,
		Input: texts,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshalling embed request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, openaiEmbedURL, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("building embed request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	res, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI embed request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("OpenAI embedding error %s: %s", res.Status, body)
	}

	var parsed openaiEmbedResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decoding embed response: %w", err)
	}

	result := make([][]float32, len(parsed.Data))
	for i, d := range parsed.Data {
		result[i] = d.Embedding
	}
	return result, nil
}
