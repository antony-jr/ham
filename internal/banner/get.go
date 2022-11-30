package banner

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji/v2"
)

func GetStartBanner() {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Print(emoji.Sprint(":fire:"), "Getting ASOP Build ", emoji.Sprint(":fire:"))
	c.Print("\n   ******************")

	fmt.Print("\n")
}

func GetRecipeNotExistsBanner() {
   	c := color.New(color.FgYellow)
	c.Print("Recipe does not exist locally, Using GIT.\n")
}

func GetFinishBanner() {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Print("Built Successfully ", emoji.Sprint(":rocket:"))
	fmt.Print("\n")
}
