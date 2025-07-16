package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type Hunk struct {
	FilePath string
	Header   string
	Body     string
}

var hunkHeaderRegex = regexp.MustCompile(`^@@\s+[-+0-9,]+\s+[-+0-9,]+\s+@@`)

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
	var currentHunk []string
	var recording bool

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "diff --git") {
			currentFile = parseFilePath(line)
			recording = false
			continue
		}

		if hunkHeaderRegex.MatchString(line) {
			if len(currentHunk) > 0 {
				hunks = append(hunks, Hunk{
					FilePath: currentFile,
					Header:   currentHunk[0],
					Body:     strings.Join(currentHunk, "\n"),
				})
			}

			currentHunk = []string{line}
			recording = true
			continue
		}

		if recording {
			currentHunk = append(currentHunk, line)
		}
	}

	if len(currentHunk) > 0 {
		hunks = append(hunks, Hunk{
			FilePath: currentFile,
			Header:   currentHunk[0],
			Body:     strings.Join(currentHunk, "\n"),
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
