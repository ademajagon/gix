package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const openaiChatURL = "https://api.openai.com/v1/chat/completions"
const openaiChatModel = "gpt-4o"
const openaiEmbedURL = "https://api.openai.com/v1/embeddings"
const openaiEmbedModel = "text-embedding-3-small"

type OpenAI struct {
	apiKey string
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{apiKey: apiKey}
}

// --- Chat completion types ---

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiChatRequest struct {
	Model    string          `json:"model"`
	Messages []openaiMessage `json:"messages"`
}

type openaiChoice struct {
	Message openaiMessage `json:"message"`
}

type openaiChatResponse struct {
	Choices []openaiChoice `json:"choices"`
}

// --- Embedding types ---

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
	prompt := CommitMessageUserPromptTemplate + diff

	payload := openaiChatRequest{
		Model: openaiChatModel,
		Messages: []openaiMessage{
			{
				Role:    "system",
				Content: CommitMessageSystemPrompt,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openaiChatURL, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("OpenAI API error (%s): %s", res.Status, string(bodyBytes))
	}

	var response openaiChatResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}

func (o *OpenAI) GetEmbeddings(texts []string) ([][]float32, error) {
	reqBody := openaiEmbedRequest{
		Model: openaiEmbedModel,
		Input: texts,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", openaiEmbedURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI embedding error %s", res.Status)
	}

	var parsed openaiEmbedResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	result := make([][]float32, len(parsed.Data))
	for i, d := range parsed.Data {
		result[i] = d.Embedding
	}

	return result, nil
}
