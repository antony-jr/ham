package cli

import (
   "os"
   "github.com/mkideal/cli"
   
   "github.com/antony-jr/ham/internal/cmd/initialize"
)

type rootT struct {
   cli.Helper
}

func Run() error {
   var root = &cli.Command {
      Name: "Hetzner Android Make (HAM)",
      Desc: "Build Android ROMs from Source with ease using Hetzner Cloud.",
      Argv: func() interface{} { return new(rootT) },
      Fn: func(ctx *cli.Context) error {
	 //argv := ctx.Argv().(*rootT)
	 return nil
      },
   }

   return cli.Root(
      root,
      cli.Tree(initialize.NewCommand()),
   ).Run(os.Args[1:])
}
