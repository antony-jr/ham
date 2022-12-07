package main

import (
	"fmt"
	"os"

	"github.com/antony-jr/ham/internal/banner"
	"github.com/antony-jr/ham/internal/cli/build_cli"
)

/*
 * Inject in build with,
 * -ldflags "-X main.AppVersion=$VERSION"
 */
var AppVersion = "Unknown"
var GitCommit = "Unknown"

func main() {
	banner.Header(AppVersion, GitCommit)
	if err := build_cli.Run(); err != nil {
		banner.Error(fmt.Sprint(err))
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		// Show usage help.
		banner.Usage()
	}
}
