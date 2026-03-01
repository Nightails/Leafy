package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Supported formats
var (
	audioFormats = []string{".mp3", ".wav", ".flac"}
	videoFormats = []string{".mp4", ".avi", ".mkv"}
)

// GetFiles scans the provided file paths and returns a list of file paths.
func GetFiles(paths []string) ([]string, error) {
	var mediaFiles []string

	for _, p := range paths {
		if p == "" {
			continue
		}

		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			mediaFiles = addFile(mediaFiles, p)
			continue
		}

		err = filepath.WalkDir(p, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			mediaFiles = addFile(mediaFiles, path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return mediaFiles, nil
}

// addFile adds a file path to the given list if it's a supported format and not already present.
func addFile(files []string, path string) []string {
	if !isSupported(path) {
		return files
	}
	if slices.Contains(files, path) {
		return files
	}
	files = append(files, path)
	return files
}

// isSupported returns true if the given file path has a supported format.
func isSupported(path string) bool {
	for _, ext := range audioFormats {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			return true
		}
	}
	for _, ext := range videoFormats {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			return true
		}
	}
	return false
}
