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

type chatClient struct {
	chatURL    string
	embedURL   string
	chatModel  string
	embedModel string
	apiKey     string
	httpClient *http.Client
}

func newChatClient(chatURL, embedURL, chatModel, embedModel, apiKey string, timeout time.Duration) *chatClient {
	return &chatClient{
		chatURL:    chatURL,
		embedURL:   embedURL,
		chatModel:  chatModel,
		embedModel: embedModel,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: timeout},
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float32       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
	Error   *chatError   `json:"error"`
}

type chatError struct {
	Message string `json:"message"`
}

type embedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embedResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func (c *chatClient) GenerateCommitMessage(diff string) (string, error) {
	payload := chatRequest{
		Model: c.chatModel,
		Messages: []chatMessage{
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

	req, err := http.NewRequest(http.MethodPost, c.chatURL, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		var apiErr chatResponse
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != nil {
			return "", fmt.Errorf("API error: %s", apiErr.Error.Message)
		}
		return "", fmt.Errorf("API error: %s", res.Status)
	}

	var response chatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return strings.TrimSpace(response.Choices[0].Message.Content), nil
}

func (c *chatClient) GetEmbeddings(texts []string) ([][]float32, error) {
	payload := embedRequest{
		Model: c.embedModel,
		Input: texts,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshalling embed request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.embedURL, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("creating embed request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embed request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("embed API error %s: %s", res.Status, body)
	}

	var parsed embedResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decoding embed response: %w", err)
	}

	result := make([][]float32, len(parsed.Data))
	for i, d := range parsed.Data {
		result[i] = d.Embedding
	}
	return result, nil
}
