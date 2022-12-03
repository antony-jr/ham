package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// Short SHA256 Sum is used since sha256 sum is simply too
// long for a server name in Hetzner.
// It's concatenation of first 7 characters and last 7 characters
// of the original sha256 sum hex.
func ServerNameFromSHA256(sum string) string {
	shortSHA256Sum := sum[:7] + sum[len(sum)-7:]
	serverName := fmt.Sprintf("build-%s", shortSHA256Sum)
	return serverName
}

// Remove all servers at hetzner project which is running
// beyond 24 hours. This is because we want to reduce cost
// at all times and we really don't want dead expesive servers
// running around wasting our money.
// So to accomplish that we simply destroy servers older than
// 24 hours or with age of 24 hours.
// We destroy servers whenever we see them or connect with
// the api.
//
// "If you are 1 day old then you are dead to me."
// -- Antony J.R
func DestroyAllDeadServers(sclient *hcloud.ServerClient) error {
	servers, err := sclient.All(
		context.Background(),
	)

	if err != nil {
		return err
	}

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}
	now := time.Now().In(loc)

	for _, server := range servers {
		diff := now.Sub(server.Created)
		hours := int(diff.Hours())

		if hours >= 24 {
			_, _, _ = sclient.DeleteWithResult(
				context.Background(),
				server,
			)
		}
	}

	return nil
}
