package initialize

import (
	"context"
	"errors"
	"fmt"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"

	"github.com/antony-jr/ham/internal/banner"
	"github.com/antony-jr/ham/internal/core"
	"github.com/antony-jr/ham/internal/helpers"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mkideal/cli"
)

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

			configFilePath, err := helpers.ConfigFilePath()
			if err != nil {
				return errors.New("Failed to Get Config File Path")
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

			pks := string(publicKey[:])
			pks = fmt.Sprint(pks[:len(pks)-2])

			config := core.NewConfiguration(
				argv.APIKey,
				fmt.Sprintf("%s= ham@antonyjr.in\n", pks),
				string(privateKeyBytes[:]),
			)

			// Add the new sshkey and check connection with the API Key
			client := hcloud.NewClient(hcloud.WithToken(config.APIKey))

			sshkeys, err := client.SSHKey.All(
				context.Background(),
			)
			if err != nil {
				return err
			}

			// Search for ham-ssh-key SSH Key,
			// if found then delete it if the user
			// forces it.
			keyExists := false
			for _, el := range sshkeys {
				if el.Name == "ham-ssh-key" {
					keyExists = true
					break
				}
			}

			if keyExists {
				if !argv.Force {
					return errors.New("SSH Key Already Exists at Server, Run with -f flag.")
				}

				// Delete the existing key first then add our new one.
				_, err := client.SSHKey.Delete(
					context.Background(),
					sshkeys[0],
				)
				if err != nil {
					return err
				}
			}

			labels := make(map[string]string)

			// Add the new key now
			_, _, err = client.SSHKey.Create(
				context.Background(),
				hcloud.SSHKeyCreateOpts{
					"ham-ssh-key",
					config.SSHPublicKey,
					labels,
				},
			)

			if err != nil {
				return err
			}

			// Create Configuration File
			exists, err := helpers.FileExists(configFilePath)
			if err != nil {
				return err
			}

			if exists && !argv.Force {
				return errors.New("Configuration Already Exists, Run with -f flag.")
			}

			err = core.WriteConfiguration(config)
			if err != nil {
				return errors.New("Cannot Write Configuration File")
			}

			banner.InitFinishBanner()
			return nil
		},
	}
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
