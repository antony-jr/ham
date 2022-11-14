package get

import (
   	"os"
	"fmt"
	"github.com/mkideal/cli"
)

type getT struct {
	cli.Helper
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
			//fmt.Println("Getting")
			args := ctx.Args()

			if len(args) != 1 {
			   return nil
			}

			fmt.Println("Getting ", args[0])

			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			   fmt.Println("Recipe does not exist locally, trying remote.")
			}

			//ctx.String("Hello, root command, I am %s\n", argv.Name)
			
			// argv := ctx.Argv().(*initT)
			// argv.name
			return nil
		},
	}
}
