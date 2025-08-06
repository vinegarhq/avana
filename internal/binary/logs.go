package binary

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func robloxLogFile() (string, error) {
	ad, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("get appdata: %w", err)
	}

	dir := filepath.Join(ad, "Roblox", "logs")
	slog.Info("Searching for log file", "dir", dir)

	// This is required due to fsnotify requiring the directory
	// to watch to exist before adding it.
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create roblox log dir: %w", err)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return "", fmt.Errorf("make fsnotify watcher: %w", err)
	}
	defer w.Close()

	if err := w.Add(dir); err != nil {
		return "", fmt.Errorf("watch roblox log dir: %w", err)
	}

	t := time.NewTimer(logTimeout)

	for {
		select {
		case <-t.C:
			return "", fmt.Errorf("roblox log file not found after %s", logTimeout)
		case e := <-w.Events:
			if e.Has(fsnotify.Create) {
				return e.Name, nil
			}
		case err := <-w.Errors:
			slog.Error("Recieved fsnotify watcher error", "error", err)
		}
	}
}