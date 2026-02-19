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

	return diff, nil
}
