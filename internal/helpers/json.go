package helpers

import (
	"encoding/json"
	"errors"
	"os"
)

func DumpJsonFile(obj map[string]string, path string) error {
	json, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.New("Cannot Write: " + path)
	} else {
		file.Write(json)
		file.Close()
	}
	return nil
}

func DumpJsonString(obj map[string]string, path string) (string, error) {
	json, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return string(json), nil
}
