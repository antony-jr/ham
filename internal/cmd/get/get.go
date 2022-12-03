package get

import (
   	"os"
	"fmt"
	"context"
	"errors"
	"strings"

	"github.com/mkideal/cli"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"golang.org/x/crypto/ssh"

	"github.com/charmbracelet/lipgloss"

	"github.com/antony-jr/ham/internal/banner"
	"github.com/antony-jr/ham/internal/core"
	"github.com/antony-jr/ham/internal/helpers"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type getT struct {
	cli.Helper
	Force  bool   `cli:"f,force" usage:"Force start a build even if the recipe was built Already."`
}

func ParseGitRemoteString(remote string) (string, string) {
   urlSlice := strings.Split(remote, "://")

   if len(urlSlice) == 2 {
      remote = urlSlice[1]
   }

   slice := strings.Split(remote, "/")

   branchSlice := strings.Split(remote, ":")
   branch := ""
   url := ""

   if len(branchSlice) == 2 {
      url = branchSlice[0]
      branch = branchSlice[1]
   } else {
      url = remote
   }

   if len(urlSlice) == 2 {
      url = fmt.Sprintf("%s://%s", urlSlice[0], url)
   }

   if len(slice) != 2 {
      // Regular Git url
      return url, branch
   }

   user := slice[0]
   userSlice := strings.Split(user, "@")

   if len(userSlice) != 2 {
      return url, branch
   }

   uname := userSlice[0]
   host := userSlice[1]

   if strings.ToLower(host) != "gh" {
      return url, branch
   }

   repo := slice[1]
   repoSlice := strings.Split(repo, ":")

   if len(repoSlice) == 2 {
      repo = repoSlice[0]
   }

   return fmt.Sprintf("https://github.com/%s/%s", uname, repo), branch
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "get",
		Desc: "Get a build of ASOP from community recipe or locally using your Hetzner Cloud",
		Text:`
Syntax: ham get [RECIPE LOCATION]

Recipe from Github:
   ham get user@gh/repo:branch
   ham get antony-jr@gh/enchilada_los18.1
   ham get antony-jr@gh/ecnhilada_los18.1:dev

Recipe from Git:
   ham get https://antonyjr.in/enchilada_los181.git

Local Recipe:
   ham get ./examples/enchilada_los18.1`,
		Argv: func() interface{} { return new(getT) },
		NumArg: func(n int) bool {
		   if n != 1 {
		      return false
		   }
		   return true
		},
		Fn: func(ctx *cli.Context) error {
		   	argv := ctx.Argv().(*getT)
			args := ctx.Args()
			if len(args) != 1 {
			   return nil
			}
			recipe_src := args[0]
			dir := recipe_src

			checkMark := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")

			banner.GetStartBanner()
			
			fmt.Printf(" %s Parsing %s...\n", checkMark, recipe_src)
			remove := false
			if _, err := os.Stat(recipe_src); os.IsNotExist(err) {
			   // Recipe is not local, so use git to clone the 
			   // the recipe requested by the user.
			   banner.GetRecipeNotExistsBanner()

			   // Parse the string
			   git_url, git_branch := ParseGitRemoteString(recipe_src)
			   orig_branch := git_branch

			   if git_branch == "" {
			      git_branch = "default"
			   }

			   fmt.Printf(" %s Git URL: %s\n", checkMark, git_url)
			   fmt.Printf(" %s Git Branch: %s\n", checkMark, git_branch) 
			   

			   uniqueTempDir, err := os.MkdirTemp(os.TempDir(), "*-ham-recipe")
			   if err != nil {
			      return err
			   }
			   dir = uniqueTempDir
			   remove = true

			   fmt.Printf(" %s Cloning Into: %s\n", checkMark, dir)

			   r, err := git.PlainClone(dir, false, &git.CloneOptions{
			      URL: git_url,
			   })

			   if err != nil {
			      _ = os.RemoveAll(dir) 
			      return err
			   }

			   if orig_branch != "" {
			      w, err := r.Worktree()
			      if err != nil {
				 _ = os.RemoveAll(dir) 	
				 return err
			      }

			      err = w.Checkout(&git.CheckoutOptions{Branch: plumbing.ReferenceName(orig_branch)})
			      if err != nil {
				 _ = os.RemoveAll(dir)	
				 return err
			      }
			   }
			}

			if remove {
			    defer os.RemoveAll(dir)
			}

			config, err := core.GetConfiguration()
			if err != nil {
			   return err
			}
			fmt.Printf(" %s Reading Configuration\n", checkMark)


			client := hcloud.NewClient(hcloud.WithToken(config.APIKey))
			fmt.Printf(" %s Connected with Hetzner Cloud API\n", checkMark)


			sshkeys, err := client.SSHKey.All(
				context.Background(),
			)
			if err != nil {
				return err
			}

			// Search for ham-ssh-key SSH Key,
			// if it does not exists then error out
			// asking the user to properly init.
			keyOk := false
			keyFingerprint, err := sshFingerprint(config.SSHPublicKey)

			// fmt.Println("SSH Key Fingerprint: ", keyFingerprint)

			if err != nil {
			   return err
			}

			var ham_labels map[string]string

			for _, el := range sshkeys {
				if el.Name == "ham-ssh-key" {
				   	// fmt.Println("Hetzner Key Fingerprint: ", el.Fingerprint)
				        if keyFingerprint == el.Fingerprint {
					   keyOk = true
					   ham_labels = el.Labels
					}	
					break
				}
			}

			if !keyOk {
			   return errors.New("HAM SSH Key not found, Please Re-Initialize.")
			}

			fmt.Printf(" %s Verified SSH Keys\n", checkMark)


			// Destroy all dead servers 
			// whenver we see them.
			// This is highly unlikely that our ham leaves dead servers
			// but this is just a precaution.
			err = helpers.DestroyAllDeadServers(&client.Server)
			if err != nil {
			   return err
			}

			// Parse recipe file for meta information 
			// and args information.
			hf, err := core.NewHAMFile(dir)
			if err != nil {
				return err
			}
			serverName := helpers.ServerNameFromSHA256(hf.SHA256Sum)

			// Search for build servers that were already started
			// if found, track that.
			servers, err := client.Server.All(
			   context.Background(),
			)
			if err != nil {
			   return err
			}

			for _, server := range servers {
			   if server.Name == serverName {
			      // Track status instead of creating a new one.
			      fmt.Printf(" %s Active Build Found\n", checkMark)
			      fmt.Printf(" %s Started Progress on Active Build\n", checkMark)
			      return nil
			   }
			}

			// Check the ham-ssh-key labels, label with the server
			// name will have the last build status like success
			// or failed.
			for key, status := range ham_labels {
			   if key == serverName && !argv.Force {
			      // Confirm with user before starting the 
			      // build again.
			      estr := fmt.Sprintf("A %s build had run before with this recipe, Run with -f flag to force build.",
			      			  status)
			      return errors.New(estr)
			   }
			   break
			}
			fmt.Printf(" %s Checked Previous Builds\n", checkMark)
			fmt.Println()

			// Create a new build server.

			// Before that we need to get variables from the user
			// such as special files, env vars required for the 
			// build from the user. This might be crucial secrets
			// so transport it with SSH to stay secure.

			// TODO: Actually create a server and ssh into that host.
			sshClient, err := GetSSHClient("[::1]:22", config.SSHPrivateKey)
			if err != nil {
			   return err
			}
			defer sshClient.Close()

			sshSession, err := GetSSHSession(sshClient)
			if err != nil {
			   return err
			}
			defer sshSession.Close()

			shell, err := GetSSHShell(sshSession)
			if err != nil {
			   return err
			}

			/* Update and Upgrade */
			_, _ = shell.Exec("export DEBIAN_FRONTEND=noninteractive")
			_, _ = shell.Exec("apt-get update -y --force-yes -qq")
			_, _ = shell.Exec("apt-get upgrade -y --force-yes -qq")

			err = runProgressTeaProgram(shell)
			if err != nil {
			   return err
			}

			return nil
		},
	}
}

func sshFingerprint(pubkey string) (string, error) {
   pubKeyBytes := []byte(pubkey)

   pk, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
   if err != nil {
      return "", err
   }

   return ssh.FingerprintLegacyMD5(pk), nil
}
