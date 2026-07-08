package cli

import (
	"errors"
	"fmt"
	"os"
)

func validatePort(port int) error {
	if port < 1024 || port > 65535 {
		return errors.New("port must be between 1024 and 65535")
	}

	return nil
}

func validateDir(dir string) error {
	info, err := os.Stat(dir)

	if err != nil {
		return fmt.Errorf("Directory error: %w", err)
	}

	if !info.IsDir() {
		return errors.New("provided path is not a directory")
	}

	return nil
}

func ValidateConfig(config Config) error {
	if err := validatePort(config.Port); err != nil {
		return err
	}

	return validateDir(config.Directory)
}
