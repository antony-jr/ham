package core

import (
	"context"
	"errors"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// Change this if needed in the future when
// Hetzner deprecates Ubuntu 20.04 LTS, or if it
// is that time of the year.
const (
	TargetImage = "ubuntu-20.04"

	// Most Stable, Reliable and Cheapest
	// Location by Hetzner
	TargetLocation = "nbg1"
)

func CreateServer(client *hcloud.Client, server *hcloud.ServerType, serverName string) (*hcloud.Server, error) {
	// Get Server Image
	serverImage, _, err := client.Image.Get(
		context.Background(),
		TargetImage,
	)
	if err != nil {
		return nil, err
	}

	// Get ham-ssh-key SSH Key
	sshKey, _, err := client.SSHKey.Get(
		context.Background(),
		"ham-ssh-key",
	)
	if err != nil {
		return nil, err
	}

	// Get Location
	location, _, err := client.Location.Get(
		context.Background(),
		TargetLocation,
	)
	if err != nil {
		return nil, err
	}

	startAfterCreate := true
	// Server Creation Options
	serverCreateOpts := hcloud.ServerCreateOpts{
		Name:             serverName,
		ServerType:       server,
		Image:            serverImage,
		SSHKeys:          []*hcloud.SSHKey{sshKey},
		Location:         location,
		StartAfterCreate: &startAfterCreate,
		Labels:           map[string]string{},
		PublicNet: &hcloud.ServerCreatePublicNet{
			EnableIPv4: false,
			EnableIPv6: true,
		},
	}

	err = serverCreateOpts.Validate()
	if err != nil {
		return nil, err
	}

	// Create Server at Hetzner
	createResult, _, err := client.Server.Create(
		context.Background(),
		serverCreateOpts,
	)

	if err != nil {
		return nil, err
	}

	// Wait till we Success or Failure
	// result from Action that is currently
	// running.
	ok := false
	errMsg := ""

	// Check Current Action First
	checkAction(createResult.Action, &ok, &errMsg)
	if !ok {
		return nil, errors.New(errMsg)
	}

	// Loop Over all Next Actions and
	// See if it's ok
	for _, action := range createResult.NextActions {
		checkAction(action, &ok, &errMsg)
		if !ok {
			return nil, errors.New(errMsg)
		}
	}

	return createResult.Server, nil
}

func checkAction(action *hcloud.Action, ok *bool, errMsg *string) {
	*ok = false
	*errMsg = ""
	for {
		if action.Status == hcloud.ActionStatusRunning {
			time.Sleep(time.Second * time.Duration(2))
			continue
		} else if action.Status == hcloud.ActionStatusSuccess {
			*ok = true
		} else if action.Status == hcloud.ActionStatusError {
			*ok = false
			*errMsg = "Server Create Action Failed (" + action.ErrorMessage + ")"
		}
		break
	}
}
