package get

import (
	"golang.org/x/crypto/ssh"
	//"golang.org/x/crypto/ssh/knownhosts"
)

type SSHShellCode int

const (
	SSH_SHELL_NO_ERROR SSHShellCode = iota
	SSH_SHELL_CANNOT_GET_CLIENT
	SSH_SHELL_CANNOT_GET_SESSION
	SSH_SHELL_CANNOT_CONNECT
	SSH_SHELL_MALFORMED_JSON
	SSH_SHELL_HAM_STATUS_ERRORED
)

type SSHShellContext struct {
	client *ssh.Client
	code   SSHShellCode
}

func GetSSHClient(host string, privKey string) (*ssh.Client, error) {
	pKey := []byte(privKey)

	var err error
	var signer ssh.Signer

	signer, err = ssh.ParsePrivateKey(pKey)
	if err != nil {
		return nil, err
	}

	/*
	   var hostkeyCallback ssh.HostKeyCallback
	   hostkeyCallback, err = knownhosts.New("known_hosts")
	   if err != nil {
	      return nil, err
	   }*/

	conf := &ssh.ClientConfig{
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	var conn *ssh.Client
	conn, err = ssh.Dial("tcp", host, conf)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetSSHSession(conn *ssh.Client) (*ssh.Session, error) {
	var session *ssh.Session
	var err error

	session, err = conn.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func GetSSHShell(client *ssh.Client) (*SSHShellContext, error) {
	return &SSHShellContext{
		client: client,
	}, nil
}

func (ctx *SSHShellContext) Exec(command string) (string, error) {
	session, err := GetSSHSession(ctx.client)
	if err != nil {
		return "", err
	}
	defer session.Close()

	out, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (ctx *SSHShellContext) SetCode(c SSHShellCode) {
	ctx.code = c
}
