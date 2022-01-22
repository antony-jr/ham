package initialize

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"

	"github.com/antony-jr/ham/internal/banner"
	"github.com/mkideal/cli"
)

type Configuration struct {
	APIKey        string
	SSHPublicKey  string
	SSHPrivateKey string
}

type initT struct {
	cli.Helper
	APIKey string `pw:"k,key" usage:"Hetzner API Key/Token for the project" prompt:"Hetzner API Key/Token"`
	Force  bool   `cli:"f,force" usage:"Overwrite configuration even if exists."`
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "init",
		Desc: "Initialize API connection to your Hetzner Project",
		Argv: func() interface{} { return new(initT) },
		Fn: func(ctx *cli.Context) error {
			argv := ctx.Argv().(*initT)

			banner.InitStartBanner()

			configFilePath, err := getConfigFilePath()
			if err != nil {
				return errors.New("Failed to Get Config File Path")
			}

			exists, err := fileExists(configFilePath)
			if err != nil {
				return err
			}

			if exists && !argv.Force {
				return errors.New("Configuration Already Exists, Run with -f flag.")
			}

			if len(argv.APIKey) < 10 {
				return errors.New("Invalid API Key")
			}

			privateKey, err := generatePrivateKey(4096)
			if err != nil {
				return errors.New("SSH Key Generation Failed")
			}

			publicKey, err := generatePublicKey(&privateKey.PublicKey)
			if err != nil {
				return errors.New("SSH Key Generation Failed")
			}

			privateKeyBytes := encodePrivateKeyToPEM(privateKey)

			config := Configuration{
				argv.APIKey,
				fmt.Sprintf("%s ham@antonyjr.in\n", string(publicKey[:len(publicKey)-2])),
				string(privateKeyBytes[:]),
			}

			json, err := json.Marshal(config)
			if err != nil {
				return errors.New("JSON Encoding Failed")
			}

			// Create Configuration File
			file, err := os.Create(configFilePath)
			if err != nil {
				return errors.New("Cannot Write Configuration File")
			} else {
				file.Write(json)
				file.Close()
			}

			banner.InitFinishBanner()
			return nil
		},
	}
}

func fileExists(FilePath string) (bool, error) {
	if _, err := os.Stat(FilePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%c.ham.json", homedir, os.PathSeparator), nil
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}
