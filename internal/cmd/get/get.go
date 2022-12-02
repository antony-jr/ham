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

	"github.com/antony-jr/ham/internal/banner"
	"github.com/antony-jr/ham/internal/core"
	"github.com/antony-jr/ham/internal/helpers"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type getT struct {
	cli.Helper
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
			args := ctx.Args()
			if len(args) != 1 {
			   return nil
			}
			recipe_src := args[0]
			dir := recipe_src

			banner.GetStartBanner()
			fmt.Println("Recipe: ", recipe_src)
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

			   fmt.Println("Git URL: ", git_url)
			   fmt.Println("Git Branch: ", git_branch) 
			   

			   uniqueTempDir, err := os.MkdirTemp(os.TempDir(), "*-ham-recipe")
			   if err != nil {
			      return err
			   }
			   dir = uniqueTempDir
			   remove = true

			   fmt.Println("Cloning Into: ", dir)

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

			client := hcloud.NewClient(hcloud.WithToken(config.APIKey))

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

			for _, el := range sshkeys {
				if el.Name == "ham-ssh-key" {
				   	// fmt.Println("Hetzner Key Fingerprint: ", el.Fingerprint)
				        if keyFingerprint == el.Fingerprint {
					   keyOk = true
					}	
					break
				}
			}

			if !keyOk {
			   return errors.New("HAM SSH Key not found, Please Re-Initialize.")
			}


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

			      return nil
			   }
			}

			// Create a new build server.

			// Before that we need to get variables from the user
			// such as special files, env vars required for the 
			// build from the user. This might be crucial secrets
			// so transport it with SSH to stay secure.

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
