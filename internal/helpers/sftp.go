package helpers

import (
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func SFTPFileExistsAtRemote(conn *ssh.Client, path string) bool {
	client, err := sftp.NewClient(conn)
	if err != nil {
		return false
	}
	defer client.Close()

	_, err = client.Lstat(path)
	if err != nil {
		return false
	}

	return true
}

func GetSFTPClient(conn *ssh.Client) (*sftp.Client, error) {
	return sftp.NewClient(conn)
}

func SFTPCopyFileToRemote(client *sftp.Client, dest string, source string) error {
	f, err := client.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer f.Close()

	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = f.ReadFrom(src)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
