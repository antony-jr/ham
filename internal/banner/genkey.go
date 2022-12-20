package banner

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji/v2"
)

func GenKeyStartBanner(c, s, l, o, ou, cn, e string, ks int) {
	in := "# Generating Android Certificates\n"
	in += "Country: **%s**\n\n"
	in += "State: **%s**\n\n"
	in += "Locality: **%s**\n\n"
	in += "Organization: **%s**\n\n"
	in += "Organization Unit: **%s**\n\n"
	in += "Common Name: **%s**\n\n"
	in += "E-Mail: %s\n\n"
	in += "Key Algorithm: **RSA %d**\n\n"

	in = fmt.Sprintf(in, c, s, l, o, ou, cn, e, ks)

	out, _ := glamour.Render(in, "dark")
	fmt.Print(out)
}

func GenKeyFinishBanner() {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Print("Certificates Generated Successfully ", emoji.Sprint(":rocket:"))
	fmt.Print("\n")
}
