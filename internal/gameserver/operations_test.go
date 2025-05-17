package gameserver

import (
	"context"
	"fmt"
	"testing"

	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

/* func TestCreateGameServerInternal(t *testing.T) {
	helper, err := connhelper.GetConnectionHelper("ssh://ajay@dockerhost.local")

	if err != nil {
		return
	}

	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
	}

	var clientOpts []client.Opt

	clientOpts = append(clientOpts,
		client.WithHTTPClient(httpClient),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
	)
	// Create a mock WorldItem
	worldItem := &models.WorldItem{
		Name: "test-world",
	}

	// Call the CreateGameServer function
	resp, err := createGameServerInternal(worldItem, clientOpts)

	// Check for errors
	if err != nil {
		t.Fatalf("createGameServerInternal returned an error: %v", err)
	}

	// Check if the response is not empty
	if resp.ID == "" {
		t.Fatal("createGameServerInternal returned an empty response ID")
	}

	t.Cleanup(func() {
		if resp.ID == "" {
			return
		}
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return
		}
		defer cli.Close()

		_ = cli.ContainerStop(context.Background(), resp.ID, container.StopOptions{})
	})
} */

func TestCreateGameServer(t *testing.T) {

	// Create a mock WorldItem
	worldItem := &models.WorldItem{
		Name: "test-world-2",
	}

	// Call the CreateGameServer function
	resp, err := CreateGameServer(worldItem)

	// Check for errors
	if err != nil {
		t.Fatalf("CreateGameServer returned an error: %v", err)
	}

	// Check if the response is not empty
	if resp.ID == "" {
		t.Fatal("CreateGameServer returned an empty response ID")
	}

	t.Cleanup(func() {
		if resp.ID == "" {
			return
		}
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return
		}
		defer cli.Close()

		e := cli.ContainerStop(context.Background(), resp.ID, container.StopOptions{})

		if e != nil {
			fmt.Println(fmt.Errorf("error stopping container: ID=%v: %v", resp.ID, e))
		} else {
			fmt.Printf("Container stopped: ID=%v\n", resp.ID)
		}

		e = cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
		if e != nil {
			fmt.Println(fmt.Errorf("error removing container: ID=%v: %v", resp.ID, e))
		} else {
			fmt.Printf("Container removed: ID=%v\n", resp.ID)
		}
	})
}
