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

const geminiBaseURL = "https://generativelanguage.googleapis.com/v1beta/models"
const geminiChatModel = "gemini-flash-latest"
const geminiEmbedModel = "gemini-embedding-001"

type Gemini struct {
	apiKey string
}

func NewGemini(apiKey string) *Gemini {
	return &Gemini{apiKey: apiKey}
}

// --- Chat types ---

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiChatRequest struct {
	SystemInstruction *geminiContent  `json:"systemInstruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}

type geminiChatResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
}

// --- Embedding types ---

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
}

func (g *Gemini) GenerateCommitMessage(diff string) (string, error) {
	url := fmt.Sprintf("%s/%s:generateContent", geminiBaseURL, geminiChatModel)

	prompt := CommitMessageUserPromptTemplate + diff

	payload := geminiChatRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{
				{Text: CommitMessageSystemPrompt},
			},
		},
		Contents: []geminiContent{
			{
				Role:  "user",
				Parts: []geminiPart{{Text: prompt}},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", g.apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("Gemini API error (%s): %s", res.Status, string(bodyBytes))
	}

	var response geminiChatResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no response from Gemini")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}

func (g *Gemini) GetEmbeddings(texts []string) ([][]float32, error) {
	url := fmt.Sprintf("%s/%s:batchEmbedContents", geminiBaseURL, geminiEmbedModel)

	modelRef := fmt.Sprintf("models/%s", geminiEmbedModel)

	requests := make([]geminiEmbedContentRequest, len(texts))
	for i, t := range texts {
		requests[i] = geminiEmbedContentRequest{
			Model: modelRef,
			Content: geminiContent{
				Parts: []geminiPart{{Text: t}},
			},
		}
	}

	payload := geminiBatchEmbedRequest{Requests: requests}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", g.apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("Gemini embedding error (%s): %s", res.Status, string(bodyBytes))
	}

	var parsed geminiBatchEmbedResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	result := make([][]float32, len(parsed.Embeddings))
	for i, e := range parsed.Embeddings {
		result[i] = e.Values
	}

	return result, nil
}
