package config

import (
	"testing"
)

func setTestHome(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
}

func TestConfig_ResolveProvider_Default(t *testing.T) {
	cfg := Config{}
	if cfg.ResolveProvider() != "openai" {
		t.Errorf("expected default provider 'openai', got %q", cfg.ResolveProvider())
	}
}

func TestConfig_ResolveProvider_Explicit(t *testing.T) {
	cases := []string{"openai", "gemini", "ollama"}
	for _, p := range cases {
		cfg := Config{Provider: p}
		if cfg.ResolveProvider() != p {
			t.Errorf("expected %q, got %q", p, cfg.ResolveProvider())
		}
	}
}

func TestConfig_APIKey_OpenAI(t *testing.T) {
	cfg := Config{Provider: "openai", OpenAIKey: "sk-test"}
	if cfg.APIKey() != "sk-test" {
		t.Errorf("expected OpenAI key, got %q", cfg.APIKey())
	}
}

func TestConfig_APIKey_Gemini(t *testing.T) {
	cfg := Config{Provider: "gemini", GeminiKey: "gemini-test"}
	if cfg.APIKey() != "gemini-test" {
		t.Errorf("expected Gemini key, got %q", cfg.APIKey())
	}
}

func TestConfig_APIKey_Ollama(t *testing.T) {
	cfg := Config{Provider: "ollama"}
	if cfg.APIKey() != "" {
		t.Errorf("expected empty key for ollama, got %q", cfg.APIKey())
	}
}

func TestConfig_APIKey_DefaultsToOpenAI(t *testing.T) {
	cfg := Config{OpenAIKey: "sk-default"}
	if cfg.APIKey() != "sk-default" {
		t.Errorf("expected OpenAI key as default, got %q", cfg.APIKey())
	}
}

func TestSaveAndLoad(t *testing.T) {
	setTestHome(t)

	cfg := Config{
		OpenAIKey:          "sk-test",
		GeminiKey:          "gemini-test",
		Provider:           "ollama",
		OllamaBaseURL:      "http://localhost:11434",
		OllamaChatModel:    "mistral",
		OllamaEmbedModel:   "nomic-embed-text",
		DisableUpdateCheck: true,
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loaded.Provider != "ollama" {
		t.Errorf("expected provider 'ollama', got %q", loaded.Provider)
	}
	if loaded.OllamaBaseURL != "http://localhost:11434" {
		t.Errorf("unexpected OllamaBaseURL: %q", loaded.OllamaBaseURL)
	}
	if loaded.OllamaChatModel != "mistral" {
		t.Errorf("unexpected OllamaChatModel: %q", loaded.OllamaChatModel)
	}
	if loaded.OllamaEmbedModel != "nomic-embed-text" {
		t.Errorf("unexpected OllamaEmbedModel: %q", loaded.OllamaEmbedModel)
	}
	if !loaded.DisableUpdateCheck {
		t.Error("expected DisableUpdateCheck to be true")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	setTestHome(t)
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing config file, got nil")
	}
}

func TestSave_CreatesDirectoryIfMissing(t *testing.T) {
	setTestHome(t)

	cfg := Config{Provider: "ollama"}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}
}

func TestOllamaFields_RoundTrip(t *testing.T) {
	setTestHome(t)

	original := Config{
		Provider:         "ollama",
		OllamaBaseURL:    "http://remote:11434",
		OllamaChatModel:  "llama3.1:8b",
		OllamaEmbedModel: "nomic-embed-text",
	}

	if err := Save(original); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if loaded.OllamaBaseURL != original.OllamaBaseURL {
		t.Errorf("OllamaBaseURL mismatch: got %q, want %q", loaded.OllamaBaseURL, original.OllamaBaseURL)
	}
	if loaded.OllamaChatModel != original.OllamaChatModel {
		t.Errorf("OllamaChatModel mismatch: got %q, want %q", loaded.OllamaChatModel, original.OllamaChatModel)
	}
	if loaded.OllamaEmbedModel != original.OllamaEmbedModel {
		t.Errorf("OllamaEmbedModel mismatch: got %q, want %q", loaded.OllamaEmbedModel, original.OllamaEmbedModel)
	}
}
