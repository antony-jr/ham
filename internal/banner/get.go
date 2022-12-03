package banner

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji/v2"
)

func GetStartBanner() {
	in := `# Get Build`

	out, _ := glamour.Render(in, "dark")
	fmt.Print(out)
}

func GetMalformedJSONBanner(serverName string) {
	in := "# Tracking Failed\n"
	in += "Cannot get builder status but the build is **still running** at **%s**, please use the following command to\n"
	in += "destory all server currently running in the project regardless if it's created by ham or try again to track\n"
	in += "the progress with ```ham get recipe```.\n"
	in += "```\n"
	in += " $ ham clean \n"
	in += "```\n"
	in += "\n\n"

	in = fmt.Sprintf(in, serverName)

	out, _ := glamour.Render(in, "auto")
	fmt.Print(out)
}

func GetRecipeBanner(name string, ver string, hash string) {
	in := "# Recipe Information\n"
	in += "**Name**: *%s* [%s]                           **SHA-256**: %s"
	in = fmt.Sprintf(in, name, ver, hash)

	out, _ := glamour.Render(in, "auto")
	fmt.Print(out)
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
