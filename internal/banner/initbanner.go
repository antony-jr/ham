package banner

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji/v2"
)

func InitStartBanner() {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Print(emoji.Sprint(":fire:"), "Initializing HAM ", emoji.Sprint(":fire:"))
	c.Print("\n   ****************")

	fmt.Print("\n")
}

func InitFinishBanner() {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Print("Initialized Successfully ", emoji.Sprint(":rocket:"))
	fmt.Print("\n")
}
