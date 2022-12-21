package genkey

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"math/big"

	"github.com/antony-jr/ham/internal/banner"
	"github.com/antony-jr/ham/internal/helpers"
	"github.com/mkideal/cli"
)

type genkeyT struct {
	cli.Helper
	Output     string `cli:"r,output" usage:"Output ZIP File."`
	Country    string `cli:"c,country" usage:"The Country where you Generate this Key."`
	State      string `cli:"s,state" usage:"The State  where you Generate this Key."`
	City       string `cli:"l,locality" usage:"The City  where you Generate this Key."`
	Org        string `cli:"o,org" usage:"The Organization where you Generate this Key."`
	OrgUnit    string `cli:"u,org-unit" usage:"The Organization Unit where you Generate this Key."`
	CommonName string `cli:"n,common-name" usage:"Common Name for the Keys."`
	Email      string `cli:"e,email" usage:"The Contact E-Mail for the Generated Keys."`
	KeySize    int    `cli:"k,key-size" usage:"The Key Size to Use in RSA Key."`
	Force      bool   `cli:"f,force" usage:"Overwrite AndroidCerts.zip file if it exists."`
}

type keyT struct {
	Certificate []byte
	PK8         []byte
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "genkey",
		Desc: "Generate Android Certificates and Compress them directly into ZIP",
		Argv: func() interface{} { return new(genkeyT) },
		Fn: func(ctx *cli.Context) error {
			argv := ctx.Argv().(*genkeyT)

			country := "US"
			state := "California"
			city := "Mountain View"
			org := "Android"
			orgUnit := "Android"
			commonName := "Android"
			email := "android@android.com"
			keySize := 2048
			destArchive := "AndroidCerts.zip"

			if len(argv.Country) != 0 {
				country = argv.Country
			}
			if len(argv.State) != 0 {
				state = argv.State
			}
			if len(argv.City) != 0 {
				city = argv.City
			}
			if len(argv.Org) != 0 {
				org = argv.Org
			}
			if len(argv.OrgUnit) != 0 {
				orgUnit = argv.OrgUnit
			}
			if len(argv.CommonName) != 0 {
				commonName = argv.CommonName
			}
			if len(argv.Email) != 0 {
				email = argv.Email
			}
			if argv.KeySize != 0 {
				keySize = argv.KeySize
			}

			// Check if KeySize is valid.
			if keySize != 2048 &&
				keySize != 4096 {
				return errors.New("Invalid RSA Key Size.")
			}

			if len(argv.Output) != 0 {
				destArchive = argv.Output
			}

			banner.GenKeyStartBanner(country,
				state,
				city,
				org,
				orgUnit,
				commonName,
				email,
				keySize)

			keyNames := []string{
				"releasekey",
				"platform",
				"shared",
				"media",
				"networkstack",
				"testkey",
			}

			keys := map[string]keyT{}

			for _, keyName := range keyNames {
				fmt.Printf("   - Generating %s...\n", keyName)

				cert, pk8, err := makeKey(country,
					state,
					city,
					org,
					orgUnit,
					commonName,
					email,
					keySize)
				if err != nil {
					return err
				}

				keys[keyName] = keyT{
					Certificate: cert,
					PK8:         pk8,
				}
			}

			fmt.Println()

			exists, err := helpers.FileExists(destArchive)
			if err != nil {
				return err
			}

			if exists {
				if !argv.Force {
					return errors.New(fmt.Sprintf("%s File Already Exists, Run with -f to Force Write.", destArchive))
				} else {
					err := os.Remove(destArchive)
					if err != nil {
						return err
					}
				}
			}

			archive, err := os.Create(destArchive)
			if err != nil {
				return err
			}
			defer archive.Close()

			zipWriter := zip.NewWriter(archive)
			defer zipWriter.Close()

			for keyName, key := range keys {
				certEntry := fmt.Sprintf("%s.x509.pem", keyName)
				keyEntry := fmt.Sprintf("%s.pk8", keyName)

				certFile, err := zipWriter.Create(certEntry)
				if err != nil {
					return err
				}

				_, err = io.Copy(certFile, bytes.NewReader(key.Certificate))
				if err != nil {
					return err
				}

				keyFile, err := zipWriter.Create(keyEntry)
				if err != nil {
					return err
				}

				_, err = io.Copy(keyFile, bytes.NewReader(key.PK8))
				if err != nil {
					return err
				}
			}

			banner.GenKeyFinishBanner()

			return nil
		},
	}
}

func makeKey(country string,
	state string,
	city string,
	org string,
	orgUnit string,
	cName string,
	email string,
	size int) ([]byte, []byte, error) {
	privateKey, err := generatePrivateKey(size)
	if err != nil {
		return nil, nil, err
	}

	oidEmailAddress := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
	subj := pkix.Name{
		CommonName:         cName,
		Country:            []string{country},
		Province:           []string{state},
		Locality:           []string{city},
		Organization:       []string{org},
		OrganizationalUnit: []string{orgUnit},
		ExtraNames: []pkix.AttributeTypeAndValue{
			{
				Type: oidEmailAddress,
				Value: asn1.RawValue{
					Tag:   asn1.TagIA5String,
					Bytes: []byte(email),
				},
			},
		},
	}

	now := time.Now()
	future := time.Now().AddDate(10, 0, 0)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 196)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := asn1.Marshal(*privateKey.Public().(*rsa.PublicKey))
	if err != nil {
		return nil, nil, err
	}

	subjectKeyId := sha1.Sum(publicKeyBytes)

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subj,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             now,
		NotAfter:              future,
		IsCA:                  true,
		BasicConstraintsValid: true,
		SubjectKeyId:          subjectKeyId[:],
		AuthorityKeyId:        subjectKeyId[:],
	}

	csrBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, privateKey.Public(), privateKey)
	if err != nil {
		return nil, nil, err
	}

	certBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Headers: nil, Bytes: csrBytes})
	pemBytes, err := encodePrivateKeyToPEM(privateKey)

	if err != nil {
		return nil, nil, err
	}

	return certBytes, pemBytes, nil
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	// Get ASN.1 DER format
	privDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return privDER, nil
}
