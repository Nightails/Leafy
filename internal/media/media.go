package media

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Supported media file formats
var (
	audioFormats = []string{".mp3", ".wav", ".flac"}
	videoFormats = []string{".mp4", ".avi", ".mkv"}
)

// GetMediaFiles scans the provided file paths and returns a list of media file paths.
func GetMediaFiles(paths []string) ([]string, error) {
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
			mediaFiles = addMediaFile(mediaFiles, p)
			continue
		}

		err = filepath.WalkDir(p, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			mediaFiles = addMediaFile(mediaFiles, path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return mediaFiles, nil
}

// addMediaFile adds a media file path to the given list if it's a supported format and not already present.
func addMediaFile(files []string, path string) []string {
	if !isSupportedFormat(path) {
		return files
	}
	if slices.Contains(files, path) {
		return files
	}
	files = append(files, path)
	return files
}

// isSupportedFormat returns true if the given file path has a supported format.
func isSupportedFormat(path string) bool {
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
