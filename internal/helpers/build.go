package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"

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
		if !strings.HasPrefix(server.Name, "build-") {
			continue
		}

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

func TryDeleteServer(sclient *hcloud.ServerClient, serverName string, maxTries int, interval int) error {
	delTries := 0
	for {
		delErr := DeleteServer(sclient, serverName)
		if delErr.Error() == "Server Not Found" {
			break
		}

		delTries++
		fmt.Println("Destroying Server Failed. Retrying... ")
		time.Sleep(time.Second * time.Duration(interval))
		if delTries > maxTries {
			return errors.New("Cannot Destroy Remote Server. " + delErr.Error())
		}
	}

	return nil
}

func DeleteServer(sclient *hcloud.ServerClient, serverName string) error {
	servers, err := sclient.All(
		context.Background(),
	)

	if err != nil {
		return err
	}

	for _, server := range servers {
		if server.Name == serverName {
			result, _, err := sclient.DeleteWithResult(
				context.Background(),
				server)

			if err != nil {
				return err
			}

			if result.Action.Status == "error" {
				return errors.New(result.Action.ErrorMessage)
			}
			return nil
		}
	}

	return errors.New("Server Not Found")
}

func GetServerAgeInHours(sclient *hcloud.ServerClient, serverName string) (int, error) {
	servers, err := sclient.All(
		context.Background(),
	)

	if err != nil {
		return -1, err
	}

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return -1, err
	}

	for _, server := range servers {
		if server.Name == serverName {
			now := time.Now().In(loc)
			diff := now.Sub(server.Created)
			hours := int(diff.Hours())

			return hours, nil
		}
	}

	return -1, errors.New("Server Not Found")
}

func UpdateSSHKeyLabel(sclient *hcloud.SSHKeyClient, sshKey *hcloud.SSHKey, key string, value string) (*hcloud.SSHKey, error) {
	lbls := sshKey.Labels
	lbls[key] = value

	newKey, _, err := sclient.Update(
		context.Background(),
		sshKey,
		hcloud.SSHKeyUpdateOpts{
			Name:   sshKey.Name,
			Labels: lbls,
		},
	)

	return newKey, err
}

func ReadVarsJsonFile(path string) (map[string]string, error) {
	result := make(map[string]string)
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(source, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
