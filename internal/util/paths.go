package util

import (
	"os"
	"path/filepath"
)

func AppDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".mockspin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func SessionsDir() (string, error) {
	app, err := AppDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(app, "sessions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}