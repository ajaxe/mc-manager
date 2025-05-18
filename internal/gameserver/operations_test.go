package gameserver

import (
	"context"
	"fmt"
	"slices"
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

	t.Cleanup(func() { cleanupContainer(resp, t) })
}

func TestGameServerIntance(t *testing.T) {
	n := "test-world-hats-mayname"
	worldItem := &models.WorldItem{
		Name: n,
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

	result, err := GameServerIntance()

	if err != nil {
		t.Fatalf("GameServerIntance returned an error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("GameServerIntance returned an empty list of names")
	}
	if slices.Contains(result, n) == false {
		t.Fatalf("GameServerIntance returned an unexpected name: got %v, want %v", result, n)
	}

	t.Cleanup(func() { cleanupContainer(resp, t) })
}

func cleanupContainer(resp container.CreateResponse, t *testing.T) {
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
		t.Logf("%v\n", fmt.Errorf("error stopping container: ID=%v: %v", resp.ID, e))
	} else {
		t.Logf("Container stopped: ID=%v\n", resp.ID)
	}

	e = cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
	if e != nil {
		t.Logf("%v\n", fmt.Errorf("error removing container: ID=%v: %v", resp.ID, e))
	} else {
		t.Logf("Container removed: ID=%v\n", resp.ID)
	}
}
