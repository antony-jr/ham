package get

import (
	"fmt"
	"github.com/mkideal/cli"
)

type initT struct {
	cli.Helper
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "get",
		Desc: "Get a build of ASOP from community recipe or locally using your Hetzner Cloud",
		Argv: func() interface{} { return new(initT) },
		Fn: func(ctx *cli.Context) error {
			fmt.Println("Getting")
			// argv := ctx.Argv().(*initT)
			// argv.name
			return nil
		},
	}
}
