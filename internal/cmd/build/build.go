package build

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/antony-jr/ham/internal/banner"
	"github.com/antony-jr/ham/internal/core"
	"github.com/mkideal/cli"
	"github.com/sevlyar/go-daemon"
)

type buildT struct {
	cli.Helper
	Sum        string `cli:"*s,sum" usage:"SHA256 Hash of the main ham.yaml file"`
	RecipePath string `cli:"*r,recipe" usage:"Recipe file path which has the ham.yaml"`
	VarsPath   string `cli:"*a,vars" usage:"JSON file path containing all required build variables prompted"`
}

type statusT struct {
	Status string
	Title  string
	Error  error
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "build",
		Desc: "Build ASOP from source from a recipe and given build vars. (*Run in Build Machine)",
		Argv: func() interface{} { return new(buildT) },
		Fn: func(ctx *cli.Context) error {
			argv := ctx.Argv().(*buildT)
			if runtime.GOOS != "linux" {
				return errors.New("OS Not Supported.")
			}

			hf, err := core.NewHAMFile(argv.RecipePath)
			if err != nil {
				return err
			}

			banner.BuildStartBanner()

			fmt.Printf("%s\n", hf.Title)
			fmt.Printf("v%s\n", hf.Version)
			fmt.Printf("SHA256 Sum: %s\n", hf.SHA256Sum)

			if hf.SHA256Sum != argv.Sum {
				return errors.New("SHA256 Mismatch, Bad File.")
			}

			//serverName := fmt.Sprintf("build-%s", hf.SHA256Sum)

			// We assume the current server name at hetzner to
			// be this and we use this assumption to destroy
			// the server when the build is done.


			dctx := &daemon.Context{
				PidFileName: "/tmp/com.github.antony-jr.ham.pid",
				PidFilePerm: 0644,
				LogFileName: "/tmp/com.github.antony-jr.ham.log",
				LogFilePerm: 0640,
				WorkDir:     "./",
				Umask:       027,
				Args: []string{"ham",
					"build",
					"-r",
					argv.RecipePath,
					"-a",
					argv.VarsPath,
					"-s",
					argv.Sum},
			}

			d, err := dctx.Reborn()
			if err != nil {
				return err
			}

			if d != nil {
				banner.BuildFinishBanner()
				return nil
			}

			// Daemon Execution
			// Actual Builder
			defer destroyCurrentServer(hf.SHA256Sum)

			// This holds the status in json,
			// the TCP server responds with this
			// status string when asked
			status := statusT{
				"Running",
				"",
				nil,
			}

			go statusServer(&status)

			terminal, err := core.NewTerminal(hf.SHA256Sum)
			if err != nil {
				return err
			}
			defer terminal.CloseTerminal()

			for index, el := range hf.Build {
				status.Status = "Building"
				status.Title = el.Title
				if status.Error != nil {
					return status.Error
				}

				err := terminal.ExecTerminal(index, el.Cmd)
				if err != nil {
					return err
				}

				err = terminal.WaitTerminal(index)
				if err != nil {
					return err
				}
			}

			status.Status = "Uploading"
			status.Title = "Build Finished, Uploading Assets"
			fmt.Println("Built Successfully.")

			// After build completes.
			// Upload the outputs to a very cheap server
			// by Hetzner, and name is serve
			// This reduces cost a lot

			// The serve server will have the name
			// assets-<sha256 sum of yaml>
			// The get command should detect this new server
			// and serve the user that and delete the cheap
			// server and complete the build entirely.

			// We need to destroy this current expensive server
			// at the end.
			// The current server should be named
			// build-<sha256 sum of yaml>

			return nil
		},
	}
}

func destroyCurrentServer(UniqueID string) {
	serverName := fmt.Sprintf("build-%s", UniqueID)
	fmt.Println("Destroying ", serverName)
}

func statusServer(state *statusT) {
	listener, err := net.Listen("tcp", "0.0.0.0:1695")
	if err != nil {
		state.Error = err
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleRequest(state, conn)
	}
}

func handleRequest(state *statusT, conn net.Conn) {
	buf := make([]byte, 1024)
	rLen, err := conn.Read(buf)
	if err != nil {
		return
	}

	request := strings.ToLower(string(buf[:rLen]))
	var resp string
	if request == "status" {
		resp = fmt.Sprintf("{ error: false, status: \"%s\", progress: \"%s\" }\n", state.Status, state.Title)
	} else if request == "quit" {
		resp = fmt.Sprintf("{ error: false, status: \"Stopping\", progress: \"Stopping\"}\n")
		defer os.Exit(0)
	} else {
		resp = fmt.Sprintf("{ error: true, message: \"Unknown command\" }\n")
	}

	conn.Write([]byte(resp))
	conn.Close()
}
