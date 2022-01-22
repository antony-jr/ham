package banner

import (
   "fmt"

   "github.com/kyokomi/emoji/v2"
   "github.com/fatih/color"
)

func Header(Version string) {
   cliName := color.New(color.FgRed).Add(color.Bold)
   cliName.Print("Ham ", emoji.Sprint(":hamster:"))

   fmt.Print("(")
   ver := color.New(color.Bold)
   ver.Print(Version)
   fmt.Print(")")

   fmt.Print(",")

   tagLine := color.New(color.FgGreen).Add(color.Bold)
   tagLine.Println(" Hetzner Android Make.")

   copyright := color.New(color.FgCyan).Add(color.Bold)
   copyright.Println("Copyright (c) 2022, D. Antony J.R.\n")
}
