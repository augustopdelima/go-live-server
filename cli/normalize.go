package cli

import (
	"os"
	"path/filepath"
	"strings"
)

func expandPath(path string) (string, error) {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return home, nil
	}

	if strings.HasPrefix(path, "~/"){
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(
			home, path[2:],
		), nil
	}

	return path, nil
}



func NormalizeConfig(config Config) (Config, error) {

	directory, err := expandPath(config.Directory)

	if err != nil {
		return  config, err
	}

	directory, err = filepath.Abs(directory)

	if err != nil {
		return  config, err
	}

	config.Directory = directory;

	return config, nil
}
