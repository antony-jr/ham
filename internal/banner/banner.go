package banner

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji/v2"
)

func Header(Version string, Commit string) {
	cliName := color.New(color.FgRed).Add(color.Bold)
	cliName.Print("Ham ", emoji.Sprint(":hamster:"))

	fmt.Print("(v")
	ver := color.New(color.Bold)
	ver.Print(Version)
	fmt.Print(" commit-")
	commit := color.New(color.FgYellow).Add(color.Bold)
	commit.Print(Commit)
	fmt.Print("),")

	tagLine := color.New(color.FgGreen).Add(color.Bold)
	tagLine.Println(" Hetzner Android Make.")

	copyright := color.New(color.FgCyan).Add(color.Bold)
	copyright.Println("Copyright (c) 2022, D. Antony J.R <antonyjr@protonmail.com>.")
	copyright.Println("The BSD 3-Clause \"New\" or \"Revised\" License.")
	fmt.Print("\n")
}

func Usage() {
	c := color.New(color.FgWhite).Add(color.Bold)
	c.Print("Usage: ")

	fmt.Println(os.Args[0], " [OPTIONS] [COMMAND]")
	fmt.Println("Run help command to get more information")
}

func Error(Message string) {
	c := color.New(color.FgRed).Add(color.Bold)
	c.Print("Fatal: ", Message)

	fmt.Print("\n")
}
