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

	fmt.Println("🔐 Stashing working directory...")
	if err := exec.Command("git", "stash", "push", "--include-untracked", "--keep-index", "-m", "gix-temp").Run(); err != nil {
		return fmt.Errorf("git stash failed: %w", err)
	}
	defer func() {
		fmt.Println("Restoring stashed changes...")
		exec.Command("git", "stash", "pop").Run()
	}()

	for i, group := range groups {
		fmt.Printf("Committing group %d: %s\n", i+1, group.Message)

		if err := exec.Command("git", "reset").Run(); err != nil {
			return fmt.Errorf("git reset failed: %w", err)
		}

		patch := JoinGroupPatch(group.Hunks)
		tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("gix_group_%d.patch", i))
		if err := os.WriteFile(tmpPath, []byte(patch), 0600); err != nil {
			return fmt.Errorf("write patch failed: %w", err)
		}

		cmd := exec.Command("git", "apply", "--cached", tmpPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git apply failed: %w", err)
		}

		commit := exec.Command("git", "commit", "-m", group.Message)
		commit.Stdout = os.Stdout
		commit.Stderr = os.Stderr
		if err := commit.Run(); err != nil {
			return fmt.Errorf("git commit failed: %w", err)
		}
	}

	return nil
}
