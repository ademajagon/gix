package semantics

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ApplyGroups(groups []HunkGroup) error {
	if len(groups) == 0 {
		return fmt.Errorf("no group to apply")
	}

	stash := exec.Command("git", "stash", "push", "--include-untracked", "--keep-index", "--quiet", "-m", "gix-temp")
	if err := stash.Run(); err != nil {
		return fmt.Errorf("git stash failed: %w", err)
	}

	defer func() {
		_ = exec.Command("git", "stash", "pop", "--quiet").Run()
	}()

	for i, group := range groups {
		fmt.Printf("[commit %d/%d] %s\n", i+1, len(groups), group.Message)

		if err := exec.Command("git", "reset").Run(); err != nil {
			return fmt.Errorf("git reset failed: %w", err)
		}

		patch := JoinGroupPatch(group.Hunks)
		tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("gix_group_%d.patch", i))
		if err := os.WriteFile(tmpPath, []byte(patch), 0600); err != nil {
			return fmt.Errorf("write patch failed: %w", err)
		}

		if err := exec.Command("git", "apply", "--cached", tmpPath).Run(); err != nil {
			return fmt.Errorf("git apply failed: %w", err)
		}

		if err := exec.Command("git", "commit", "-m", group.Message).Run(); err != nil {
			return fmt.Errorf("git commit failed: %w", err)
		}
	}

	return nil
}
