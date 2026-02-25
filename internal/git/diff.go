package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// IsGitRepo checks whether the current working dir is inside a git work tree.
func IsGitRepo() bool {
	return exec.Command("git", "rev-parse", "--is-inside-work-tree").Run() == nil
}

// HasStagedChanges checks whether there are any staged changes.
func HasStagedChanges() (bool, error) {
	out, err := exec.Command("git", "diff", "--cached", "--name-only").Output()
	if err != nil {
		return false, fmt.Errorf("git diff --cached --name-only: %w", err)
	}

	return strings.TrimSpace(string(out)) != "", nil
}

const maxDiffBytes = 3_000

func GetStagedDiff() (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("git", "diff", "--cached", "--unified=3")
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git diff --cached: %w", err)
	}

	diff := strings.TrimSpace(buf.String())
	if diff == "" {
		return "", fmt.Errorf("no staged changes - `git add <files>` to stage changes")
	}

	return truncateDiff(diff, maxDiffBytes), nil
}

func truncateDiff(diff string, maxBytes int) string {
	if len(diff) <= maxBytes {
		return diff
	}

	if idx := lastIndexWithin(diff, "\ndiff --git", maxBytes); idx > 0 {
		truncated := strings.TrimSpace(diff[:idx])
		total := countFiles(diff)
		kept := countFiles(truncated)
		return fmt.Sprintf("%s\n\n[diff truncated: showing %d of %d file(s)]", truncated, kept, total)
	}

	if idx := lastIndexWithin(diff, "\n@@", maxBytes); idx > 0 {
		truncated := strings.TrimSpace(diff[:idx])
		return fmt.Sprintf("%s\n\n[diff truncated: showing partial diff of first file]", truncated)
	}

	return diff[:maxBytes] + "\n\n[diff truncated]"
}

func lastIndexWithin(s, marker string, maxBytes int) int {
	idx := -1
	start := 0
	for {
		i := strings.Index(s[start:], marker)
		if i < 0 {
			break
		}
		abs := start + i
		if abs >= maxBytes {
			break
		}
		idx = abs
		start = abs + len(marker)
	}
	return idx
}

func countFiles(diff string) int {
	return strings.Count(diff, "\ndiff --git") + 1
}
