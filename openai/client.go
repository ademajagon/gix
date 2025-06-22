package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
	prompt := "Based on the following Git diff, generate a concise commit message:\n\n" + diff

	payload := RequestPayload{
		Model: defaultModel,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant that writes clean, conventional commit messages.",
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
		return "", errors.New("OpenAI API error: " + string(bodyBytes))
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
