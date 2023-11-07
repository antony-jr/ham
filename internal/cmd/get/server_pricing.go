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
	// CCX33
	// vCPU: 8
	// Memory: 32 GB
	// Disk: 240 GB
	// Additional Volume: 400 GB for Lineage Build
	return GrossServerPriceForServerType(client, "ccx33" /*"cpx51"*/)
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
			if server.ServerType.Cores >= 8 &&
				server.ServerType.Memory >= 32.0 &&
				server.ServerType.Disk >= 240 {
				return amount, server.ServerType, nil
			}
		}
	}

	return 0.0, nil, errors.New("Cannot Find Suitable Server")
}
