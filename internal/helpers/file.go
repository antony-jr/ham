package helpers

import (
	"errors"
	"fmt"
	"os"
)

func FileExists(FilePath string) (bool, error) {
	if _, err := os.Stat(FilePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func ConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%c.ham.json", homedir, os.PathSeparator), nil
}
