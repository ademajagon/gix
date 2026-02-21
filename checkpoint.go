package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ademajagon/gix/checkpoint"
	"github.com/ademajagon/gix/config"
)

func init() {
	checkpointResult = make(chan *checkpoint.CheckResponse, 1)
}

var checkpointResult chan *checkpoint.CheckResponse

// runCheckpoint is called as a goroutine from main() before cobra runs
func runCheckpoint(currentVersion string) {
	cfg, err := config.Load()
	if err != nil {
		checkpointResult <- nil
		return
	}

	if cfg.DisableUpdateCheck {
		log.Printf("[INFO] Update check disabled via config.")
		checkpointResult <- nil
		return
	}

	configDir, cacheDir, err := checkpointDirs()
	if err != nil {
		log.Printf("[ERR] Checkpoint setup: %s", err)
		checkpointResult <- nil
		return
	}

	resp, err := checkpoint.Check(&checkpoint.CheckParams{
		Product:       "gix",
		Version:       currentVersion,
		SignatureFile: filepath.Join(configDir, "gix_id"),
		CacheFile:     filepath.Join(cacheDir, "checkpoint_cache"),
	})

	if err != nil {
		log.Printf("[ERR] Checkpoint: %s", err)
		checkpointResult <- nil
		return
	}

	checkpointResult <- resp
}

func showCheckpointResult() {
	resp := <-checkpointResult
	if resp == nil || !resp.Outdated {
		return
	}

	fmt.Fprintf(os.Stderr,
		"\n"+
			"─────────────────────────────────────────────────────\n"+
			"  Your version of gix is out of date!\n"+
			"\n"+
			"  Current version: %s\n"+
			"  Latest version:  %s\n"+
			"\n"+
			"  Update: %s\n"+
			"─────────────────────────────────────────────────────\n",
		currentVersionString(),
		resp.CurrentVersion,
		resp.CurrentDownloadURL,
	)
}

func currentVersionString() string {
	return storedVersion
}

var storedVersion string

func checkpointDirs() (configDir, cacheDir string, err error) {
	cfgBase, err := os.UserConfigDir()
	if err != nil {
		return "", "", fmt.Errorf("user config dir: %w", err)
	}

	cchBase, err := os.UserCacheDir()
	if err != nil {
		return "", "", fmt.Errorf("user cache dir: %w", err)
	}

	return filepath.Join(cfgBase, "gix"), filepath.Join(cchBase, "gix"), nil
}
