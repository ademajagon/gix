package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

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
