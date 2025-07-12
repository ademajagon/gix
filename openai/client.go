package openai

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

const apiURL = "https://api.openai.com/v1/chat/completions"
const defaultModel = "gpt-4o"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestPayload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ResponsePayload struct {
	Choices []Choice `json:"choices"`
}

func GenerateCommitMessage(apiKey string, diff string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("missing OpenAI API key.\nRun:\n  gix config set-key <your-api-key>")
	}

	prompt := "Write a single-line conventional commit message that describes the following Git diff. Only return the commit message. Do not include explanations, newlines, or formatting beyond the message itself. Diff:\n\n" + diff

	payload := RequestPayload{
		Model: defaultModel,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a concise assistant that only returns a one-line, conventional commit message. No explanations, markdown, or commentary."},
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

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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

	var response ResponsePayload
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}

func GroupDiffIntoCommits(apiKey string, diff string) ([]string, error) {
	if apiKey == "" {
		return nil, errors.New("missing OpenAI API key")
	}

	prompt := `You are helping a developer organize their Git changes into structured commits.

The diff contains changes across multiple files. Some files may work together as part of the same feature, even if they're in different packages (e.g., "openai/client.go" + "utils/utils.go").

Your task is to:
- Group changes into a **minimal number of logical commits**
- Each group should reflect a **single feature, bugfix, or task**
- Prefer **grouping files together** if they contribute to the same functional goal

Output a JSON array like this:
[
  {
    "title": "feat(ai): implement diff grouping",
    "files": ["openai/client.go", "utils/utils.go"]
  },
  ...
]

DO NOT use one commit per file unless they are truly unrelated.
DO NOT return markdown or code fences.
DO NOT include commentary or explanations — only raw JSON.

Diff:
` + diff

	payload := RequestPayload{
		Model: defaultModel,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are an expert software engineer helping organize Git history.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("OpenAI error: %s", string(body))
	}

	var response ResponsePayload
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	raw := response.Choices[0].Message.Content
	lines := strings.Split(strings.TrimSpace(raw), "\n")

	var result []string
	for _, line := range lines {
		clean := strings.TrimSpace(strings.TrimPrefix(line, "-"))
		if clean != "" && !strings.HasPrefix(clean, "```") {
			result = append(result, clean)
		}
	}

	return result, nil
}
