package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/antony-jr/ham/internal/helpers"
)

type Configuration struct {
	APIKey        string
	SSHPublicKey  string
	SSHPrivateKey string
}

func NewConfiguration(Key string, SSHPubKey string, SSHPrivKey string) Configuration {
	return Configuration{
		Key,
		SSHPubKey,
		SSHPrivKey,
	}
}

func WriteConfiguration(config Configuration) error {
	json, err := json.Marshal(config)
	if err != nil {
		return err
	}

	path, err := helpers.ConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.New("Cannot Write Configuration File")
	} else {
		file.Write(json)
		file.Close()
	}

	return nil
}

func GetConfiguration() (Configuration, error) {
	config := Configuration{}
	path, err := helpers.ConfigFilePath()
	if err != nil {
		return config, err
	}

	source, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(source, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
