package get

import (
	//"fmt"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func GrossServerPriceForServerWithHighestPerformance(client *hcloud.Client) (float64, *hcloud.ServerType, error) {
	// CPX51
	// vCPU: 16
	// Memory: 32 GB
	// Disk: 360 GB
	return GrossServerPriceForServerType(client, "cpx51" /*"cx11"*/)
}

func GrossServerPriceForServerType(client *hcloud.Client, serverType string) (float64, *hcloud.ServerType, error) {
	// Note: Adjust this if wanted in the future.
	// DE Nuremberg
	// is the most cheapest and stable server for hetzner
	// also, since this datacenter is in Germany we get to
	// enjoy GDPR.
	const (
		TargetLocation = "nbg1"
	)

	pricing, _, err := client.Pricing.Get(context.Background())
	if err != nil {
		return 0.0, nil, err
	}

	for _, server := range pricing.ServerTypes {
		// CCX are Dedicated Servers which cost a lot
		// of money, it's just better avoid them, maybe
		// in the future we might have an option to use it
		// for rich users xD.
		if !strings.HasPrefix(strings.ToLower(server.ServerType.Name), "ccx") {
			// fmt.Printf("Server Name: %s\n", server.ServerType.Name)
			// fmt.Printf("Server Cores: %d\n", server.ServerType.Cores)
			// fmt.Printf("Server Memory: %f\n", server.ServerType.Memory)
			// fmt.Printf("Server Disk: %d\n", server.ServerType.Disk)

			var hourlyPrice hcloud.Price
			priceAvail := false
			// fmt.Printf("Locations: ")
			for _, entry := range server.Pricings {
				// fmt.Printf("%s ", entry.Location.Name)
				if strings.ToLower(entry.Location.Name) == TargetLocation {
					hourlyPrice = entry.Hourly
					priceAvail = true
					break
				}
			}
			// fmt.Printf("\n\n")

			if !priceAvail {
				continue
			}

			amount, err := strconv.ParseFloat(hourlyPrice.Gross, 64)
			if err != nil {
				return 0.0, nil, errors.New("Invalid Price Given")
			}

			if len(serverType) != 0 {
				if server.ServerType.Name == serverType {
					return amount, server.ServerType, nil
				}
			} else {
				// For Some Reason, These all returns Zero
				// Maybe bug in the upstream library?
				// TODO: Look into this.
				if server.ServerType.Cores >= 16 &&
					server.ServerType.Memory >= 32.0 &&
					server.ServerType.Disk >= 320 {
					return amount, server.ServerType, nil
				}
			}
		}
	}

	return 0.0, nil, errors.New("Cannot Find Suitable Server")
}
