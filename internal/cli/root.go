package cli

import (
	"github.com/mkideal/cli"
	"os"

	"github.com/antony-jr/ham/internal/cmd/initialize"
	"github.com/antony-jr/ham/internal/cmd/build"
)

type rootT struct {
	cli.Helper
	Version bool `cli:"v,version" usage:"show version information"`
}

func Run() error {
	var root = &cli.Command{
		Name: "Hetzner Android Make (HAM)",
		Argv: func() interface{} { return new(rootT) },
		Fn: func(ctx *cli.Context) error {
			argv := ctx.Argv().(*rootT)
			if argv.Version {
				return nil
			}
			return nil
		},
	}

	var help = cli.HelpCommand("Show help")

	return cli.Root(
		root,
		cli.Tree(help),
		cli.Tree(initialize.NewCommand()),
		cli.Tree(build.NewCommand()),
	).Run(os.Args[1:])
}
