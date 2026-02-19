package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Hunk represents a single diff hunk within a file.
// Used by `gix split` (beta).
type Hunk struct {
	FilePath string
	Header   string
	Body     string
}

func ParseHunks() ([]Hunk, error) {
	cmd := exec.Command("git", "diff", "--cached", "--unified=3")
	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git diff --cached: %w", err)
	}

	return parseHunksFromDiff(buf.String())
}

// parseHunksFromDiff parses raw diff text into hunks.
func parseHunksFromDiff(diff string) ([]Hunk, error) {
	scanner := bufio.NewScanner(strings.NewReader(diff))

	var hunks []Hunk
	var currentFile string
	var fileHeader []string
	var hunkLines []string
	var hunkHeader string

	flush := func() {
		if hunkHeader == "" || len(hunkLines) == 0 {
			return
		}
		body := strings.Join(append(fileHeader, hunkLines...), "\n")
		hunks = append(hunks, Hunk{
			FilePath: currentFile,
			Header:   hunkHeader,
			Body:     body,
		})
		hunkLines = nil
		hunkHeader = ""
	}

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "diff --git "):
			flush()
			currentFile = parseFilePath(line)
			fileHeader = []string{line}
			hunkHeader = ""

		case strings.HasPrefix(line, "index ") ||
			strings.HasPrefix(line, "--- ") ||
			strings.HasPrefix(line, "+++ ") ||
			strings.HasPrefix(line, "new file") ||
			strings.HasPrefix(line, "deleted file") ||
			strings.HasPrefix(line, "old mode") ||
			strings.HasPrefix(line, "new mode"):
			fileHeader = append(fileHeader, line)

		case strings.HasPrefix(line, "@@ "):
			flush()
			hunkHeader = line
			hunkLines = []string{line}

		default:
			if hunkHeader != "" {
				hunkLines = append(hunkLines, line)
			}
		}
	}

	flush()

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning diff output: %w", err)
	}

	return hunks, nil
}

func parseFilePath(diffLine string) string {
	parts := strings.Fields(diffLine)
	if len(parts) < 4 {
		return ""
	}
	right := parts[3] // "b/foo/bar.go"
	return strings.TrimPrefix(right, "b/")
}
