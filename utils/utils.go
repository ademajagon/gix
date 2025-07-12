package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func EditInEditor(init string) string {
	tmp := filepath.Join(os.TempDir(), "gix_commit_message.txt")
	_ = os.WriteFile(tmp, []byte(init), 0600)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	fmt.Println("Editing commit message in:", editor)

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
