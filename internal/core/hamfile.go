package core

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/antony-jr/ham/internal/helpers"
	"gopkg.in/yaml.v3"
)

type HAMFile struct {
	Title     string `yaml:"title"`
	Version   string `yaml:"version"`
	SHA256Sum string
	Args      []struct {
		ID   string `yaml:"id"`
		Name string `yaml:"name"`
		Desc string `yaml:"desc"`
		Type string `yaml:"type"`
	}

	Build []struct {
		Title string `yaml:"title"`
		Cmd   string `yaml:"cmd"`
	}

	Output []struct {
		Source string `yaml:"source"`
	}
}

func NewHAMFile(RecipePath string) (HAMFile, error) {

	hf := HAMFile{}

	fp := fmt.Sprintf("%s%cham.yaml", RecipePath, os.PathSeparator)
	exists, err := helpers.FileExists(fp)
	if err != nil {
		return hf, err
	}

	if !exists {
		fp = fmt.Sprintf("%s%cham.yml", RecipePath, os.PathSeparator)
		exists, err = helpers.FileExists(fp)
		if err != nil {
			return hf, err
		}

		if !exists {
			return hf, errors.New("YAML File Not Found")
		}
	}

	source, err := ioutil.ReadFile(fp)
	if err != nil {
		return hf, err
	}

	hasher := sha256.New()
	hasher.Write(source)
	hash := hasher.Sum(nil)
	hf.SHA256Sum = fmt.Sprintf("%x", hash)

	err = yaml.Unmarshal(source, &hf)
	if err != nil {
		return hf, err
	}

	return hf, nil
}
