package split

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ApplyGroups stashes untracked changes, creates one commit per group then pops the stash
func ApplyGroups(groups []HunkGroup) error {
	if len(groups) == 0 {
		return fmt.Errorf("no groups to apply")
	}

	stash := exec.Command("git", "stash", "push",
		"--include-untracked",
		"--keep-index",
		"--quiet",
		"-m", "gix-split-temp",
	)
	stash.Stderr = os.Stderr
	if err := stash.Run(); err != nil {
		return fmt.Errorf("git stash: %w", err)
	}

	success := false
	defer func() {
		if !success {
			_ = exec.Command("git", "stash", "pop", "--quiet").Run()
		}
	}()

	for i, group := range groups {
		fmt.Printf("[%d/%d] %s\n", i+1, len(groups), group.Message)

		if err := exec.Command("git", "reset", "--quiet").Run(); err != nil {
			return fmt.Errorf("git reset (group %d): %w", i+1, err)
		}

		patch := joinPatch(group.Hunks)
		tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("gix_split_%d.patch", i))
		if err := os.WriteFile(tmpPath, []byte(patch), 0o600); err != nil {
			return fmt.Errorf("writing patch file (group %d): %w", i+1, err)
		}
		defer os.Remove(tmpPath)

		applyCmd := exec.Command("git", "apply", "--cached", tmpPath)
		applyCmd.Stderr = os.Stderr
		if err := applyCmd.Run(); err != nil {
			return fmt.Errorf("git apply (group %d): %w", i+1, err)
		}

		commitCmd := exec.Command("git", "commit", "-m", group.Message)
		commitCmd.Stdout = os.Stdout
		commitCmd.Stderr = os.Stderr
		if err := commitCmd.Run(); err != nil {
			return fmt.Errorf("git commit (group %d): %w", i+1, err)
		}
	}

	success = true

	if err := exec.Command("git", "stash", "pop", "--quiet").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: git stash pop failed: %v\n", err)
	}

	return nil
}
