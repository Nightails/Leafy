package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	audio = []string{".mp3", ".wav", ".flac"}
	video = []string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm"}
)

var formats = append(audio, video...)

func formatFileSize(bytes int64) string {
	const (
		kb = 1000
		mb = 1000 * kb
		gb = 1000 * mb
		tb = 1000 * gb
	)

	size := float64(bytes)

	switch {
	case bytes >= tb:
		return fmt.Sprintf("%.2f TB", size/tb)
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", size/gb)
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", size/mb)
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", size/kb)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func expandHome(path string) (string, error) {
	if path == "~" {
		return os.UserHomeDir()
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}
