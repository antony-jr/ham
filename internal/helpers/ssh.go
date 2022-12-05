package helpers

import "golang.org/x/crypto/ssh"

func GetSSHFingerprint(pubkey string) (string, error) {
	pubKeyBytes := []byte(pubkey)

	pk, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	if err != nil {
		return "", err
	}

	return ssh.FingerprintLegacyMD5(pk), nil
}
