package build

import (
	"fmt"
	"github.com/mkideal/cli"
	"net"
	"time"
)

type buildHaltT struct {
	cli.Helper
}

// Get clean output with:
// ham build-status | cat |  grep -a Status | cut -c 10-
func NewHaltCommand() *cli.Command {
	return &cli.Command{
		Name: "build-halt",
		Desc: "Halt or Stop Build that is currently running in the Build Machine (*Run in Build Machine) (Private)",
		Argv: func() interface{} { return new(buildStatusT) },
		Fn: func(ctx *cli.Context) error {
			_ = ctx.Argv().(*buildStatusT)

			conn, err := net.Dial("tcp", "127.0.0.1:1695")
			if err != nil {
				return err
			}

			err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				return err
			}

			// Request Status from Server
			_, err = conn.Write([]byte("quit"))
			if err != nil {
				return err
			}

			recvBuf := make([]byte, 1024)
			_, err = conn.Read(recvBuf[:])
			if err != nil {
				return err
			}

			fmt.Println("Status: ", string(recvBuf))
			return nil
		},
	}
}
