package initialize

import (
   "fmt"
   "github.com/mkideal/cli"
)

type initT struct {
   cli.Helper
}

func NewCommand() *cli.Command {
   return &cli.Command {
      Name: "init",
      Desc: "Initialize API connection to your Hetzner Project",
      Argv: func() interface {} { return new(initT) },
      Fn: func(ctx *cli.Context) error {
	 fmt.Println("Initializing")
	 // argv := ctx.Argv().(*initT)
	 // argv.name
	 return nil;
      },
   }
}
