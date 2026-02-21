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
	geminiBaseURL    = "https://generativelanguage.googleapis.com/v1beta/models"
	geminiChatModel  = "gemini-2.5-flash-latest"
	geminiEmbedModel = "gemini-embedding-001"
)

type Gemini struct {
	apiKey     string
	httpClient *http.Client
}

func NewGemini(apiKey string) *Gemini {
	return NewGeminiWithClient(apiKey, &http.Client{Timeout: 20 * time.Second})
}

func NewGeminiWithClient(apiKey string, client *http.Client) *Gemini {
	return &Gemini{apiKey: apiKey, httpClient: client}
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiChatRequest struct {
	SystemInstruction *geminiContent   `json:"systemInstruction,omitempty"`
	Contents          []geminiContent  `json:"contents"`
	GenerationConfig  *geminiGenConfig `json:"generationConfig,omitempty"`
}

type geminiGenConfig struct {
	Temperature     float32 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}

type geminiChatResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
	Error      *geminiAPIError   `json:"error,omitempty"`
}

type geminiAPIError struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Embedding types
type geminiEmbedContentRequest struct {
	Model   string        `json:"model"`
	Content geminiContent `json:"content"`
}

type geminiBatchEmbedRequest struct {
	Requests []geminiEmbedContentRequest `json:"requests"`
}

type geminiEmbedding struct {
	Values []float32 `json:"values"`
}

type geminiBatchEmbedResponse struct {
	Embeddings []geminiEmbedding `json:"embeddings"`
	Error      *geminiAPIError   `json:"error,omitempty"`
}

func (g *Gemini) GenerateCommitMessage(diff string) (string, error) {
	url := fmt.Sprintf("%s/%s:generateContent?key=%s", geminiBaseURL, geminiChatModel, g.apiKey)

	payload := geminiChatRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: CommitMessageSystem}},
		},
		Contents: []geminiContent{
			{Role: "user", Parts: []geminiPart{{Text: CommitMessageUser + diff}}},
		},
		GenerationConfig: &geminiGenConfig{
			Temperature:     0,
			MaxOutputTokens: 128,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshalling request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Gemini request: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		var apiErr geminiChatResponse
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != nil {
			return "", fmt.Errorf("Gemini API error: %s", apiErr.Error.Message)
		}
		return "", fmt.Errorf("Gemini API error %s", res.Status)
	}

	var response geminiChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("decoding Gemini response: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("Gemini returned no content")
	}

	return strings.TrimSpace(response.Candidates[0].Content.Parts[0].Text), nil
}

func (g *Gemini) GetEmbeddings(texts []string) ([][]float32, error) {
	url := fmt.Sprintf("%s/%s:batchEmbedContents?key=%s", geminiBaseURL, geminiEmbedModel, g.apiKey)
	modelRef := fmt.Sprintf("models/%s", geminiEmbedModel)

	requests := make([]geminiEmbedContentRequest, len(texts))
	for i, t := range texts {
		requests[i] = geminiEmbedContentRequest{
			Model:   modelRef,
			Content: geminiContent{Parts: []geminiPart{{Text: t}}},
		}
	}

	data, err := json.Marshal(geminiBatchEmbedRequest{Requests: requests})
	if err != nil {
		return nil, fmt.Errorf("marshalling embed request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("building embed request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Gemini embed request: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		var apiErr geminiBatchEmbedResponse
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != nil {
			return nil, fmt.Errorf("Gemini embedding error: %s", apiErr.Error.Message)
		}
		return nil, fmt.Errorf("Gemini embedding error %s", res.Status)
	}

	var parsed geminiBatchEmbedResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decoding embed response: %w", err)
	}

	result := make([][]float32, len(parsed.Embeddings))
	for i, e := range parsed.Embeddings {
		result[i] = e.Values
	}
	return result, nil
}
