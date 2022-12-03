package get

import (
	"bufio"
	"io"

	"golang.org/x/crypto/ssh"
	//"golang.org/x/crypto/ssh/knownhosts"
)

type SSHShellContext struct {
   session *ssh.Session
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
      User: "root",
      HostKeyCallback: ssh.InsecureIgnoreHostKey(),
      Auth: []ssh.AuthMethod{
	 ssh.PublicKeys(signer),
      },
   }

   var conn *ssh.Client
   conn, err = ssh.Dial("tcp6", host, conf)
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

func GetSSHShell(session *ssh.Session) (*SSHShellContext, error) {
   var err error
   session.Shell()

   _, err = session.StdinPipe()
   if err != nil {
      return nil, err
   }

   _, err = session.StdoutPipe()
   if err != nil {
      return nil, err
   }

   _, err = session.StderrPipe()
   if err != nil {
      return nil, err
   }

   return &SSHShellContext {
      session: session,
   }, nil
}

func (ctx *SSHShellContext) Exec(command string) (string, error) {
   var stdin io.WriteCloser
   var stdout io.Reader
   var err error

   stdin, err = ctx.session.StdinPipe()
   if err != nil {
      return "", err
   }

   stdout, err = ctx.session.StdoutPipe()
   if err != nil {
      return "", err
   }

   _, err = ctx.session.StderrPipe()
   if err != nil {
      return "", err
   }

   raw_command := []byte(command)

   _, err = stdin.Write(raw_command)
   if err != nil {
      return "", err
   }
  
   scanner := bufio.NewScanner(stdout)

   if tkn := scanner.Scan(); tkn {
      rcv := scanner.Bytes()

      raw := make([]byte, len(rcv))
      copy(raw, rcv)

      return string(raw), nil
   } else {
      if scanner.Err() != nil {
	 return "", scanner.Err()
      } else {
	 return "", nil
      }
   }
}


