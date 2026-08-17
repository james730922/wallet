package conf

import (
	"os"
	"path/filepath"
)

// configDir locates conf.d from the current directory or one of its parents.
// This keeps local commands and package tests independent of their working directory.
func configDir(fileName string) string {
	dir, err := os.Getwd()
	if err != nil {
		return "conf.d"
	}

	for {
		candidate := filepath.Join(dir, "conf.d")
		if _, err := os.Stat(filepath.Join(candidate, fileName)); err == nil {
			return candidate
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "conf.d"
}
