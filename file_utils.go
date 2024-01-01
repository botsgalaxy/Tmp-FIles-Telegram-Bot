package main

import (
	"os"
	"path/filepath"
)

func cleanDirectory(directoryPath string) error {
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		entryPath := filepath.Join(directoryPath, entry.Name())

		if entry.IsDir() {
			if err := cleanDirectory(entryPath); err != nil {
				return err
			}
		} else {
			if err := os.Remove(entryPath); err != nil {
				return err
			}
		}
	}

	if err := os.Remove(directoryPath); err != nil {
		return err
	}

	return nil
}
