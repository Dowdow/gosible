package utils

import (
	"path/filepath"
)

var configDir string

func SetConfigDir(configFilePath string) {
	configDir = filepath.Dir(configFilePath)
	configDir, _ = filepath.Abs(configDir)
	configDir = filepath.Clean(configDir)
}

func ResolvePath(targetPath string) string {
	if filepath.IsAbs(targetPath) {
		return filepath.Clean(targetPath)
	}

	return filepath.Clean(filepath.Join(configDir, targetPath))
}
