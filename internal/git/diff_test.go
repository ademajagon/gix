package git

import (
	"strings"
	"testing"
)

func TestTruncateDiff_NoTruncationNeeded(t *testing.T) {
	diff := "diff --git a/foo.go b/foo.go\n@@ -1,3 +1,4 @@\n+some change\n context\n"
	result := truncateDiff(diff, 10_000)
	if result != diff {
		t.Errorf("expected diff unchanged, got %q", result)
	}
}

func TestTruncateDiff_ExactLimit(t *testing.T) {
	diff := "diff --git a/foo.go b/foo.go\n+change"
	result := truncateDiff(diff, len(diff))
	if result != diff {
		t.Errorf("expected diff unchanged at exact limit, got %q", result)
	}
}

func TestTruncateDiff_CutsAtFileBoundary(t *testing.T) {
	file1 := "diff --git a/foo.go b/foo.go\n@@ -1,3 +1,4 @@\n+added line\n context\n"
	file2 := "diff --git a/bar.go b/bar.go\n@@ -1,3 +1,4 @@\n+another change\n context\n"
	diff := file1 + file2

	result := truncateDiff(diff, len(file1)+5)

	if strings.Contains(result, "bar.go") {
		t.Error("expected bar.go to be truncated out")
	}
	if !strings.Contains(result, "foo.go") {
		t.Error("expected foo.go to be present")
	}
	if !strings.Contains(result, "showing 1 of 2 file(s)") {
		t.Errorf("expected truncation notice, got: %s", result)
	}
}

func TestTruncateDiff_CutsAtFileBoundary_ThreeFiles(t *testing.T) {
	file1 := "diff --git a/a.go b/a.go\n@@ -1 +1 @@\n+change a\n"
	file2 := "diff --git a/b.go b/b.go\n@@ -1 +1 @@\n+change b\n"
	file3 := "diff --git a/c.go b/c.go\n@@ -1 +1 @@\n+change c\n"
	diff := file1 + file2 + file3

	result := truncateDiff(diff, len(file1)+len(file2)+5)

	if strings.Contains(result, "c.go") {
		t.Error("expected c.go to be truncated out")
	}
	if !strings.Contains(result, "showing 2 of 3 file(s)") {
		t.Errorf("expected '2 of 3' truncation notice, got: %s", result)
	}
}

func TestTruncateDiff_CutsAtHunkBoundary(t *testing.T) {
	header := "diff --git a/big.go b/big.go\nindex abc..def 100644\n--- a/big.go\n+++ b/big.go\n"
	hunk1 := "@@ -1,3 +1,4 @@\n+line added\n context line\n another line\n"
	hunk2 := "@@ -10,3 +11,4 @@\n+another hunk\n context\n"
	diff := header + hunk1 + hunk2

	limit := len(header) + len(hunk1) + 5
	result := truncateDiff(diff, limit)

	if strings.Contains(result, "another hunk") {
		t.Error("expected second hunk to be truncated out")
	}
	if !strings.Contains(result, "line added") {
		t.Error("expected first hunk to be present")
	}
	if !strings.Contains(result, "showing partial diff of first file") {
		t.Errorf("expected partial truncation notice, got: %s", result)
	}
}

func TestTruncateDiff_HardCutFallback(t *testing.T) {
	// no file or hunk boundaries — falls back to hard cut
	diff := strings.Repeat("x", 5_000)
	result := truncateDiff(diff, 100)

	if len(result) > 200 {
		t.Errorf("expected hard cut result to be short, got length %d", len(result))
	}
	if !strings.Contains(result, "[diff truncated]") {
		t.Errorf("expected hard cut notice, got: %s", result)
	}
}

func TestCountFiles(t *testing.T) {
	cases := []struct {
		name string
		diff string
		want int
	}{
		{
			name: "single file",
			diff: "diff --git a/foo.go b/foo.go\n+change",
			want: 1,
		},
		{
			name: "two files",
			diff: "diff --git a/foo.go b/foo.go\n+change\ndiff --git a/bar.go b/bar.go\n+change",
			want: 2,
		},
		{
			name: "three files",
			diff: "diff --git a/a.go b/a.go\ndiff --git a/b.go b/b.go\ndiff --git a/c.go b/c.go",
			want: 3,
		},
		{
			name: "empty diff",
			diff: "",
			want: 1, // strings.Count returns 0, +1 = 1
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := countFiles(c.diff)
			if got != c.want {
				t.Errorf("countFiles() = %d, want %d", got, c.want)
			}
		})
	}
}
