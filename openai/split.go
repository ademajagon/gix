package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CommitSplit struct {
	Message string `json:"message"`
	Patch   string `json:"patch"`
}

func SuggestSplitCommits(apiKey string, diff string) ([]CommitSplit, error) {
	prompt := `
You are an AI developer tool.
Given this Git diff, split it into smaller semantic chunks, each with:
- a relevant diff patch
- a single-line conventional commit message

Return the result as a JSON array like:
[
  {"message": "feat(parser): add support for trailing commas", "patch": "<diff hunk>"},
  {"message": "fix(cli): fix bug in argument parsing", "patch": "<diff hunk>"}
]

Here is the full diff:

` + diff

	payload := RequestPayload{
		Model: defaultModel,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You split diffs into meaningful commits and return structured JSON with commit messages and patches.",
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

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var response ResponsePayload
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from OpenAI")
	}

	fmt.Println(response)

	var splits []CommitSplit
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &splits); err != nil {
		return nil, fmt.Errorf("failed to parse split JSON: %w", err)
	}

	return splits, nil
}
