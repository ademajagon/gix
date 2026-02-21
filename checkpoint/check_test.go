package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsNewer(t *testing.T) {
	cases := []struct {
		latest, current string
		want            bool
	}{
		{"v1.2.3", "v1.2.3", false},
		{"v1.2.4", "v1.2.3", true},
		{"v1.3.0", "v1.2.9", true},
		{"v2.0.0", "v1.9.9", true},
		{"v1.2.2", "v1.2.3", false},
		{"v1.0.0", "v2.0.0", false},
		{"1.2.4", "1.2.3", true},
		{"1.2.3", "1.2.3", false},
	}

	for _, c := range cases {
		got := isNewer(c.latest, c.current)
		if got != c.want {
			t.Errorf("isNewer(%q, %q) == %v, want %v", c.latest, c.current, got, c.want)
		}
	}
}

func TestParseSemver(t *testing.T) {
	cases := []struct {
		input string
		want  [3]int
	}{
		{"1.2.3", [3]int{1, 2, 3}},
		{"0.0.1", [3]int{0, 0, 1}},
		{"10.20.30", [3]int{10, 20, 30}},
		{"invalid", [3]int{0, 0, 0}},
	}
	for _, c := range cases {
		if got := parseSemver(c.input); got != c.want {
			t.Errorf("parseSemver(%q) == %v, want %v", c.input, got, c.want)
		}
	}
}

func TestReadOrCreateSignature_CreateNew(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "gix_id")

	sig1 := readOrCreateSignature(path)
	if sig1 == "" {
		t.Fatal("expected a non empty signature to be generated")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("signature file not created: %v", err)
	}
	if string(data) != sig1 {
		t.Errorf("file content %q != signature %q", string(data), sig1)
	}
}

func TestReadOrCreateSignature_ReuseExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "gix_id")

	sig1 := readOrCreateSignature(path)
	sig2 := readOrCreateSignature(path)

	if sig1 != sig2 {
		t.Errorf("signature changed between reads: %q vs %q", sig1, sig2)
	}
}

func TestReadOrCreateSignature_EmptyPath(t *testing.T) {
	sig := readOrCreateSignature("")
	if sig != "" {
		t.Errorf("expected empty string for empty path, got %q", sig)
	}
}

func TestGenerateUUID_Format(t *testing.T) {
	uuid := generateUUID()
	if len(uuid) != 36 {
		t.Errorf("expected UUID length 36, got %d: %q", len(uuid), uuid)
	}

	parts := make([]int, 0)
	for _, seg := range []string{uuid[0:8], uuid[9:13], uuid[14:18], uuid[19:23], uuid[24:36]} {
		parts = append(parts, len(seg))
	}
	expected := []int{8, 4, 4, 4, 12}
	for i, p := range parts {
		if p != expected[i] {
			t.Errorf("UUID segment %d: got length %d, want %d", i, p, expected[i])
		}
	}
}

func TestCache_WriteAndRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "checkpoint_cache")

	resp := &CheckResponse{
		Product:        "gix",
		CurrentVersion: "v1.3.0",
		Outdated:       true,
	}

	writeCache(path, resp)

	got := readCache(path, DefaultCacheDuration)
	if got == nil {
		t.Fatal("expected cached response, got nil")
	}
	if got.CurrentVersion != resp.CurrentVersion {
		t.Errorf("got version %q, want %q", got.CurrentVersion, resp.CurrentVersion)
	}
	if got.Outdated != resp.Outdated {
		t.Errorf("got Outdated=%v, want %v", got.Outdated, resp.Outdated)
	}
}

func TestCache_Expired(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "checkpoint_cache")

	resp := &CheckResponse{CurrentVersion: "v1.0.0"}
	writeCache(path, resp)

	got := readCache(path, 1*time.Nanosecond)
	if got != nil {
		t.Errorf("expected nil for expired cache, got %+v", got)
	}
}

func TestCache_Missing(t *testing.T) {
	got := readCache("/tmp/gix_does_not_exist_xyz/cache", DefaultCacheDuration)
	if got != nil {
		t.Errorf("expected nil for missing cache file, got %+v", got)
	}
}

func TestCache_EmptyPath(t *testing.T) {
	writeCache("", &CheckResponse{})
	got := readCache("", DefaultCacheDuration)
	if got != nil {
		t.Errorf("expected nil for empty cache path")
	}
}

func TestCheck_DisabledByEnvVar(t *testing.T) {
	t.Setenv("GIX_CHECKPOINT_DISABLE", "1")

	resp, err := Check(&CheckParams{
		Product: "gix",
		Version: "v1.0.0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Errorf("expected nil response when disabled, got %+v", resp)
	}
}

func TestCheck_SkippedForDevBuild(t *testing.T) {
	resp, err := Check(&CheckParams{
		Product: "gix",
		Version: "dev",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Errorf("expected nil for dev build, got %+v", resp)
	}
}

func TestCheck_SkippedForCIBuild(t *testing.T) {
	resp, err := Check(&CheckParams{
		Product: "gix",
		Version: "ci",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Errorf("expected nil for CI build, got %+v", resp)
	}
}
