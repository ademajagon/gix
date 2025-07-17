package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const embedURL = "https://api.openai.com/v1/embeddings"
const embedModel = "text-embedding-3-small"

type embedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embedResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func GetEmbeddings(apiKey string, texts []string) ([][]float32, error) {
	reqBody := embedRequest{
		Model: embedModel,
		Input: texts,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", embedURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI embedding error %s", res.Status)
	}

	var parsed embedResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	result := make([][]float32, len(parsed.Data))
	for i, d := range parsed.Data {
		result[i] = d.Embedding
	}

	return result, nil
}
