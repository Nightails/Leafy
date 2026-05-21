// Package config provides configuration management for the application.
// It manages a user's configs for ease-of-use and consistency in multiple sessions.
package config

import (
	"fmt"
	"os"
)

type Config struct {
	LogPath string
}

func Default() *Config {
	userHome, _ := os.UserHomeDir()
	return &Config{
		LogPath: fmt.Sprintf("%s/.local/state/leafy/logs/leafy.log", userHome),
	}
}
