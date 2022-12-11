package get

import (
	"fmt"
	"strings"
	"time"
	"io"

	"golang.org/x/crypto/ssh"
)

func TailRemoteStdout(host string, privKey string, sum string) chan string {
	output := make(chan string)
	go func() {
		for {
			command := fmt.Sprintf("tail -f -n 5 /tmp/%s.ham.stdout \n", sum)
			client, err := GetSSHClient(host, privKey)
			if err != nil {
				time.Sleep(time.Second * time.Duration(5))
				continue
			}
			defer client.Close()

			shell, err := GetSSHShell(client)
			if err != nil {
				time.Sleep(time.Second * time.Duration(5))
				continue
			}

			out, _ := shell.Exec(fmt.Sprintf("ls /tmp/ | grep \"%s.ham.stdout\"", sum))
			if !strings.Contains(out, "ham.stdout") {
				time.Sleep(time.Second * time.Duration(20))
				continue
			}

			session, err := client.NewSession()
			if err != nil {
				time.Sleep(time.Second * time.Duration(5))
				continue
			}
			defer session.Close()

			stdout, err := session.StdoutPipe()
			if err != nil {
				time.Sleep(time.Second * time.Duration(5))
				continue
			}

			stdin, err := session.StdinPipe()
			if err != nil {
				time.Sleep(time.Second * time.Duration(5))
				continue
			}

			modes := ssh.TerminalModes{
				ssh.ECHO:          1,     // enable echoing
				ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
				ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
			}

			err = session.RequestPty("xterm-256color", 500, 400, modes)
			if err != nil {
				time.Sleep(time.Second * time.Duration(5))
				continue
			}

			err = session.Shell()
			if err != nil {
				time.Sleep(time.Second * time.Duration(10))
				continue
			}

			stdin.Write([]byte(command))

			for {
				buf := make([]byte, 500)
				_, err = io.ReadAtLeast(stdout, buf, 500)
				if err != nil {
					break
				}
				output <- string(buf)
			}
			time.Sleep(time.Second * time.Duration(10))
		}
	}()

	return output
}
