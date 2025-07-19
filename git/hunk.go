package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Hunk struct {
	FilePath string
	Header   string
	Body     string
}

func ParseHunks() ([]Hunk, error) {
	cmd := exec.Command("git", "diff", "--cached", "--unified=3")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run git diff: %w", err)
	}

	scanner := bufio.NewScanner(&stdout)

	var hunks []Hunk
	var currentFile string
	var currentLines []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "diff --git") {
			// Flush previous diff
			if len(currentLines) > 0 && currentFile != "" {
				hunks = append(hunks, Hunk{
					FilePath: currentFile,
					Body:     strings.Join(currentLines, "\n"),
				})
			}
			currentLines = []string{line}
			currentFile = parseFilePath(line)
		} else {
			currentLines = append(currentLines, line)
		}
	}

	if len(currentLines) > 0 && currentFile != "" {
		hunks = append(hunks, Hunk{
			FilePath: currentFile,
			Body:     strings.Join(currentLines, "\n"),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hunks, nil
}

func parseFilePath(line string) string {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return ""
	}

	right := parts[3]
	if strings.HasPrefix(right, "b/") {
		return right[2:]
	}

	return right
}
