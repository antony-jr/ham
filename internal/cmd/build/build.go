package build

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/antony-jr/ham/internal/banner"
	"github.com/antony-jr/ham/internal/core"
	"github.com/antony-jr/ham/internal/helpers"

	"github.com/hetznercloud/hcloud-go/hcloud"

	"github.com/mkideal/cli"
	"github.com/sevlyar/go-daemon"
)

type buildT struct {
	cli.Helper
	Sum        string `cli:"*s,sum" usage:"SHA256 Hash of the main ham.yaml file"`
	RecipePath string `cli:"*r,recipe" usage:"Recipe file path which has the ham.yaml"`
	VarsPath   string `cli:"*a,vars" usage:"JSON file path containing all required build variables prompted"`
	KeepServer bool   `cli:"k,keep-server" usage:"Don't Destroy the Remote Server on any error."`
}

type statusT struct {
	Quit       bool
	Status     string
	Title      string
	Error      error
	Percentage int
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "build",
		Desc: "Build ASOP from source from a recipe and given build vars. (*Run in Build Machine) (Private)",
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

			// We assume the current server name at hetzner to
			// be this and we use this assumption to destroy
			// the server when the build is done.
			serverName := helpers.ServerNameFromSHA256(hf.SHA256Sum)
			fmt.Printf("Build Server: %s\n", serverName)

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

			// This holds the status in json,
			// the TCP server responds with this
			// status string when asked
			status := statusT{
				false,
				"Running",
				"",
				nil,
				0,
			}

			go statusServer(&status)

			config, err := core.GetConfiguration()
			if err != nil {
				return checkErrorStatus(&status, err)
			}
			client := hcloud.NewClient(hcloud.WithToken(config.APIKey))

			// Destroy server
			// on close.
			if !argv.KeepServer {
				defer destroyCurrentServer(&client.Server, hf.SHA256Sum)
			}

			vars, err := helpers.ReadVarsJsonFile(argv.VarsPath)
			if err != nil {
				return checkErrorStatus(&status, err)
			}

			// Get ham ssh key
			hamSSHKey, _, err := client.SSHKey.Get(
				context.Background(),
				"ham-ssh-key",
			)
			if err != nil {
				return checkErrorStatus(&status, err)
			}

			// Set Label to Indicate Progress of
			// this build.
			hamSSHKey, err = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "inprogress")
			if err != nil {
				return checkErrorStatus(&status, err)
			}

			// Install Dependencies for LineageOS build/AOSP
			{
				term, err := core.NewTerminal(hf.SHA256Sum + "-prebuild")
				if err != nil {
					hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
					return checkErrorStatus(&status, err)
				}

				// We are tracking stable LTS release of Ubuntu
				// Ubuntu 20.04 (Focal)
				deps := []string{
					"bc",
					"bison",
					"build-essential",
					"ccache",
					"curl",
					"flex",
					"g++-multilib",
					"gcc-multilib",
					"git",
					"gnupg",
					"gperf",
					"imagemagick",
					"lib32ncurses5-dev",
					"lib32readline-dev",
					"lib32z1-dev",
					"libelf-dev",
					"liblz4-tool",
					"libncurses5",
					"libncurses5-dev",
					"libsdl1.2-dev",
					"libssl-dev",
					"libxml2",
					"libxml2-utils",
					"lzop",
					"pngcrush",
					"rsync",
					"schedtool",
					"squashfs-tools",
					"xsltproc",
					"zip",
					"zlib1g-dev",
					"android-sdk-platform-tools",
				}

				// The user can also install their own deps
				// from the yaml file too. This just makes life
				// so much easier when building a striaght forward
				// build from lineage.
				dep_install_command := fmt.Sprintf("apt install -y -qq %s",
					strings.Join(deps, " "))

				commands := []string{
					"export DEBIAN_FRONTEND=noninteractive",
					"apt update -y -qq",
					"apt upgrade -y -qq",
					dep_install_command,
					"curl https://storage.googleapis.com/git-repo-downloads/repo > /usr/bin/repo",
					"chmod a+x /usr/bin/repo",
					"git config --global user.email \"ham@antonyjr.in\"",
					"git config --global user.name \"Hetzner Android Make\"",
					"echo 'export USE_CCACHE=1' >> ~/.bashrc",
					"echo 'export USE_CCACHE=1' >> ~/.profile",
					"echo 'export CCACHE_EXEC=/usr/bin/ccache' >> ~/.bashrc",
					"echo 'export CCACHE_EXEC=/usr/bin/ccache' >> ~/.profile",
					"ccache -M 50G",
					"ccache -o compression=true",
				}

				status.Status = "Installing Dependencies"
				status.Title = "Installing Dependencies"

				for indx, com := range commands {
					if status.Quit {
						hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
						time.Sleep(time.Minute * time.Duration(1))
						return errors.New("User Quit the Build")
					}

					err := term.ExecTerminal(indx, com)
					if err != nil {
						hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
						return checkErrorStatus(&status, errors.New("Prebuild Failed ("+err.Error()+")"))
					}

					err = term.WaitTerminal(indx)
					if err != nil {
						hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
						return checkErrorStatus(&status, errors.New("Prebuild Failed ("+err.Error()+")"))
					}

				}

				// Set Variables for the Build Given
				// by the User.
				count := 0
				for varName, varValue := range vars {
					varName = strings.ToUpper(varName)
					varName = strings.ReplaceAll(varName, " ", "_")
					varName = strings.ReplaceAll(varName, "-", "_")
					varValue = strings.ReplaceAll(varValue, "\"", "\\\"")
					cmd := fmt.Sprintf("echo 'export %s=\"%s\"' >> ~/.bashrc", varName, varValue)
					err := term.ExecTerminal(count, cmd)
					if err != nil {
						hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
						return checkErrorStatus(&status, errors.New("Prebuild Failed (Parsing Vars)"))
					}

					err = term.WaitTerminal(count)
					if err != nil {
						hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
						return checkErrorStatus(&status, errors.New("Prebuild Failed (Parsing Vars)"))
					}

					count++
					cmd = fmt.Sprintf("echo 'export %s=\"%s\"' >> ~/.profile", varName, varValue)
					err = term.ExecTerminal(count, cmd)
					if err != nil {
						hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
						return checkErrorStatus(&status, errors.New("Prebuild Failed (Parsing Vars)"))
					}

					err = term.WaitTerminal(count)
					if err != nil {
						hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
						return checkErrorStatus(&status, errors.New("Prebuild Failed (Parsing Vars)"))
					}

					count++
				}

				term.CloseTerminal()
			}

			// Start Executing Recipe Commands.
			terminal, err := core.NewTerminal(hf.SHA256Sum)
			if err != nil {
				hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
				return checkErrorStatus(&status, err)
			}
			defer terminal.CloseTerminal()

			// Change directory to /ham-build
			err = terminal.ExecTerminal(-1, "mkdir -p /ham-build; cd /ham-build")
			if err != nil {
				hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
				return checkErrorStatus(&status, errors.New("Cannot Change to /ham-build Directory"))
			}
			err = terminal.WaitTerminal(-1)
			if err != nil {
				hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
				return checkErrorStatus(&status, errors.New("Cannot Change to /ham-build Directory"))
			}

			buildLen := len(hf.Build)
			for index, el := range hf.Build {
				if status.Quit {
					hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
					time.Sleep(time.Minute * time.Duration(1))
					return errors.New("User Quit the Build")
				}

				status.Status = "Building"
				status.Title = el.Title
				if status.Error != nil {
					hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
					return checkErrorStatus(&status, status.Error)
				}

				err := checkErrorStatus(&status, terminal.ExecTerminal(index, el.Cmd))
				if err != nil {
					hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
					return err
				}

				err = checkErrorStatus(&status, terminal.WaitTerminal(index))
				if err != nil {
					hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
					return err
				}

				// Avoid Premature Close When Tracking
				percent := int((float32(index) / float32(buildLen))) * 100.00
				if percent >= 1.0 {
					percent = percent - 1.0
				}
				status.Percentage = percent
			}

			status.Percentage = 99
			status.Status = "Finished"
			status.Title = "Build Finished"
			fmt.Println("Built Successfully.")
			fmt.Println("Running Post Build Script... ")

			status.Status = "Post Build"
			status.Title = "Running Post Build"

			pbTerminal, err := core.NewTerminal(hf.SHA256Sum + "-postbuild")
			if err != nil {
				hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
				return checkErrorStatus(&status, err)
			}
			defer pbTerminal.CloseTerminal()

			// Change directory to /ham-build
			err = pbTerminal.ExecTerminal(-1, "cd /ham-build")
			if err != nil {
				hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
				return checkErrorStatus(&status, errors.New("Cannot Change to /ham-build Directory"))
			}
			err = pbTerminal.WaitTerminal(-1)
			if err != nil {
				hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
				return checkErrorStatus(&status, errors.New("Cannot Change to /ham-build Directory"))
			}

			for index, cmd := range hf.PostBuild {
				if status.Quit {
					hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
					time.Sleep(time.Minute * time.Duration(1))
					return errors.New("User Quit the Build")
				}

				err := pbTerminal.ExecTerminal(index, cmd)
				if err != nil {
					hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
					return checkErrorStatus(&status, errors.New("Postbuild Failed ("+err.Error()+")"))
				}

				err = pbTerminal.WaitTerminal(index)
				if err != nil {
					hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "failed")
					return checkErrorStatus(&status, errors.New("Postbuild Failed ("+err.Error()+")"))
				}
			}

			hamSSHKey, _ = helpers.UpdateSSHKeyLabel(&client.SSHKey, hamSSHKey, serverName, "successful")
			status.Percentage = 100
			status.Status = "Finished"
			status.Title = "Completed"

			fmt.Println("Finished Build")

			// Give Some Time for Clients to Fetch this Status
			time.Sleep(time.Minute * time.Duration(1))
			return nil
		},
	}
}

func checkErrorStatus(state *statusT, err error) error {
	// Set Build to Error
	// We will wait for 2 mins before we exit setting
	// the status of the build at hetzner labels.
	if err == nil {
		return nil
	}

	state.Status = "Build Failed"
	state.Title = err.Error()
	state.Error = err
	state.Percentage = 100

	time.Sleep(time.Minute * time.Duration(2))
	return err
}

func destroyCurrentServer(sclient *hcloud.ServerClient, UniqueID string) {
	serverName := helpers.ServerNameFromSHA256(UniqueID)
	fmt.Println("Destroying ", serverName)

	helpers.TryDeleteServer(sclient, serverName, 20, 5)
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
	if state.Error != nil {
		resp = fmt.Sprintf("{ error: true, message: \"%s\" }\n",
			state.Error)
	} else if request == "status" {
		resp = fmt.Sprintf("{ error: false, status: \"%s\", progress: \"%s\", percentage: %d }\n",
			state.Status, state.Title, state.Percentage)
	} else if request == "quit" {
		resp = fmt.Sprintf("{ error: false, status: \"Stopping\", progress: \"Stopping\", percentage: %d }\n",
			state.Percentage)
		state.Status = "Stopping Build"
		state.Title = "Stopping Build"
		state.Quit = true
	} else {
		resp = fmt.Sprintf("{ error: true, message: \"Unknown command\" }\n")
	}

	conn.Write([]byte(resp))
	conn.Close()
}
