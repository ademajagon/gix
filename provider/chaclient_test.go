package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestChatClient(t *testing.T, chatHandler, embedHandler http.HandlerFunc) (*chatClient, *httptest.Server, *httptest.Server) {
	t.Helper()
	chatSrv := httptest.NewServer(chatHandler)
	embedSrv := httptest.NewServer(embedHandler)

	c := newChatClient(
		chatSrv.URL,
		embedSrv.URL,
		"test-chat-model",
		"test-embed-model",
		"test-api-key",
		5*time.Second,
	)
	return c, chatSrv, embedSrv
}

func TestChatClient_GenerateCommitMessage_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Error("missing or incorrect Authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("missing Content-Type header")
		}

		resp := chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Role: "assistant", Content: "feat(auth): add login support"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	c, chatSrv, embedSrv := newTestChatClient(t, handler, nil)
	defer chatSrv.Close()
	defer embedSrv.Close()

	msg, err := c.GenerateCommitMessage("diff --git a/auth.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "feat(auth): add login support" {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestChatClient_GenerateCommitMessage_TrimsWhitespace(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Content: "  feat: add feature  \n"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	c, chatSrv, embedSrv := newTestChatClient(t, handler, nil)
	defer chatSrv.Close()
	defer embedSrv.Close()

	msg, err := c.GenerateCommitMessage("some diff")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "feat: add feature" {
		t.Errorf("expected trimmed message, got %q", msg)
	}
}

func TestChatClient_GenerateCommitMessage_NoChoices(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(chatResponse{Choices: []chatChoice{}})
	}

	c, chatSrv, embedSrv := newTestChatClient(t, handler, nil)
	defer chatSrv.Close()
	defer embedSrv.Close()

	_, err := c.GenerateCommitMessage("some diff")
	if err == nil {
		t.Fatal("expected error for empty choices, got nil")
	}
}

func TestChatClient_GenerateCommitMessage_APIError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(chatResponse{
			Error: &chatError{Message: "invalid api key"},
		})
	}

	c, chatSrv, embedSrv := newTestChatClient(t, handler, nil)
	defer chatSrv.Close()
	defer embedSrv.Close()

	_, err := c.GenerateCommitMessage("some diff")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestChatClient_GenerateCommitMessage_NoAPIKey(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("expected no Authorization header for keyless client")
		}
		resp := chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Content: "chore: update deps"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	chatSrv := httptest.NewServer(http.HandlerFunc(handler))
	embedSrv := httptest.NewServer(http.NotFoundHandler())
	defer chatSrv.Close()
	defer embedSrv.Close()

	c := newChatClient(chatSrv.URL, embedSrv.URL, "llama3.1", "nomic-embed-text", "", 5*time.Second)

	msg, err := c.GenerateCommitMessage("diff")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "chore: update deps" {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestChatClient_GetEmbeddings_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := embedResponse{
			Data: []struct {
				Embedding []float32 `json:"embedding"`
			}{
				{Embedding: []float32{0.1, 0.2, 0.3}},
				{Embedding: []float32{0.4, 0.5, 0.6}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}

	c, chatSrv, embedSrv := newTestChatClient(t, nil, handler)
	defer chatSrv.Close()
	defer embedSrv.Close()

	result, err := c.GetEmbeddings([]string{"text one", "text two"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 embeddings, got %d", len(result))
	}
	if result[0][0] != 0.1 {
		t.Errorf("unexpected first embedding value: %v", result[0][0])
	}
}

func TestChatClient_GetEmbeddings_APIError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	c, chatSrv, embedSrv := newTestChatClient(t, nil, handler)
	defer chatSrv.Close()
	defer embedSrv.Close()

	_, err := c.GetEmbeddings([]string{"text"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestChatClient_SendsCorrectModels(t *testing.T) {
	var receivedChatModel, receivedEmbedModel string

	chatHandler := func(w http.ResponseWriter, r *http.Request) {
		var req chatRequest
		json.NewDecoder(r.Body).Decode(&req)
		receivedChatModel = req.Model
		json.NewEncoder(w).Encode(chatResponse{
			Choices: []chatChoice{{Message: chatMessage{Content: "fix: something"}}},
		})
	}

	embedHandler := func(w http.ResponseWriter, r *http.Request) {
		var req embedRequest
		json.NewDecoder(r.Body).Decode(&req)
		receivedEmbedModel = req.Model
		json.NewEncoder(w).Encode(embedResponse{
			Data: []struct {
				Embedding []float32 `json:"embedding"`
			}{{Embedding: []float32{0.1}}},
		})
	}

	c, chatSrv, embedSrv := newTestChatClient(t, chatHandler, embedHandler)
	defer chatSrv.Close()
	defer embedSrv.Close()

	c.GenerateCommitMessage("diff")
	if receivedChatModel != "test-chat-model" {
		t.Errorf("expected chat model %q, got %q", "test-chat-model", receivedChatModel)
	}

	c.GetEmbeddings([]string{"text"})
	if receivedEmbedModel != "test-embed-model" {
		t.Errorf("expected embed model %q, got %q", "test-embed-model", receivedEmbedModel)
	}
}
