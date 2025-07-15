package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Hunk struct {
	FilePath string
	Content  string
}

func ExtractHunks(diff string) []Hunk {
	var hunks []Hunk
	lines := strings.Split(diff, "\n")

	var currentFile string
	var currentPatch []string
	capturing := false

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			if capturing && len(currentPatch) > 0 && currentFile != "" {
				hunks = append(hunks, Hunk{
					FilePath: currentFile,
					Content:  strings.Join(currentPatch, "\n"),
				})
			}

			currentPatch = []string{line}
			currentFile = parseFilePath(line)
			capturing = true
		} else if capturing {
			currentPatch = append(currentPatch, line)
		}
	}

	if capturing && len(currentPatch) > 0 && currentFile != "" {
		hunks = append(hunks, Hunk{
			FilePath: currentFile,
			Content:  strings.Join(currentPatch, "\n"),
		})
	}

	return hunks
}

func parseFilePath(diffLine string) string {
	parts := strings.Fields(diffLine)
	if len(parts) >= 3 {
		right := parts[3]
		return strings.TrimPrefix(right, "b/")
	}
	return ""
}

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

func HasStagedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	out, err := cmd.Output()

	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(out)) != "", nil
}

func GetStagedDiff() (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("git", "diff", "--cached", "--unified=3")
	cmd.Stdout = &stdout
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	diff := strings.TrimSpace(stdout.String())

	if diff == "" {
		return "", fmt.Errorf("no staged changes to commit.\nUse `git add <files>` to stage changes first")
	}

	return diff, nil
}

func GetFullDiff() (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("git", "diff", "--unified=3")
	cmd.Stdout = &stdout
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	diff := strings.TrimSpace(stdout.String())

	if diff == "" {
		return "", fmt.Errorf("no unstaged changes found")
	}

	return diff, nil
}
