package banner

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji/v2"
)

func BuildStartBanner() {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Print(emoji.Sprint(":fire:"), "Initializing Build ", emoji.Sprint(":fire:"))
	fmt.Print("\n")
}

func BuildFinishBanner() {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Print("Started Build Successfully ", emoji.Sprint(":rocket:"))
	fmt.Print("\n")
}
