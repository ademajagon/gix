package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EditInEditor opens the init text in the users $EDITOR (fallback: nano)
func EditInEditor(init string) string {
	tmp := filepath.Join(os.TempDir(), "gix_commit_message.txt")
	if err := os.WriteFile(tmp, []byte(init), 0o600); err != nil {
		return init
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, tmp)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	edited, err := os.ReadFile(tmp)
	if err != nil {
		return init
	}
	return strings.TrimSpace(string(edited))
}
