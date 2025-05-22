package gameserver

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/ajaxe/mc-manager/internal/config"
	"github.com/ajaxe/mc-manager/internal/models"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCreateGameServer(t *testing.T) {

	// Create a mock WorldItem
	worldItem := &models.WorldItem{
		Name: newWorldName("TestCreateGameServer"),
	}

	sut := newSut()

	// Call the CreateGameServer function
	resp, err := sut.Create(worldItem)

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
	n := newWorldName("TestGameServerIntance")
	worldItem := &models.WorldItem{
		Name: n,
	}

	sut := newSut()

	// Call the CreateGameServer function
	resp, err := sut.Create(worldItem)

	// Check for errors
	if err != nil {
		t.Fatalf("CreateGameServer returned an error: %v", err)
	}

	// Check if the response is not empty
	if resp.ID == "" {
		t.Fatal("CreateGameServer returned an empty response ID")
	}

	result, err := sut.Intances()

	if err != nil {
		t.Fatalf("GameServerIntance returned an error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("GameServerIntance returned an empty list of names")
	}
	if slices.Contains(result, ToContainerName(n)) == false {
		t.Fatalf("GameServerIntance returned an unexpected name: got %v, want %v", result, n)
	}

	t.Cleanup(func() { cleanupContainer(resp, t) })
}

func TestGameServerDetail(t *testing.T) {
	id := bson.NewObjectID()
	n := newWorldName("TestGameServerDetail")
	worldItem := &models.WorldItem{
		ID:   id.Hex(),
		Name: n,
	}

	sut := newSut()

	// Call the CreateGameServer function
	resp, err := sut.Create(worldItem)

	// Check for errors
	if err != nil {
		t.Fatalf("CreateGameServer returned an error: %v", err)
	}

	// Check if the response is not empty
	if resp.ID == "" {
		t.Fatal("CreateGameServer returned an empty response ID")
	}

	result, err := sut.Details()

	if err != nil {
		t.Fatalf("GameServerDetails returned an error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("GameServerDetails returned an empty list of names")
	}
	if result[0].Name != ToContainerName(n) {
		t.Fatalf("GameServerDetails returned an unexpected name: got %v, want %v", result[0].Name, n)
	}
	if result[0].WorldID != id.Hex() {
		t.Fatalf("GameServerDetails returned an unexpected WorlID: got %v, want %v", result[0].WorldID, id.Hex())
	}

	t.Cleanup(func() { cleanupContainer(resp, t) })
}

func TestStopGameServer(t *testing.T) {
	// Create a mock WorldItem
	worldItem := &models.WorldItem{
		Name: newWorldName("TestStopGameServer"),
	}

	sut := newSut()

	// Call the CreateGameServer function
	resp, err := sut.Create(worldItem)

	// Check for errors
	if err != nil {
		t.Fatalf("CreateGameServer returned an error: %v", err)
	}

	// Check if the response is not empty
	if resp.ID == "" {
		t.Fatal("CreateGameServer returned an empty response ID")
	}

	err = sut.Stop(worldItem)

	if err != nil {
		t.Fatalf("StopGameServer returned an error: %v", err)
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

func newSut() *GameServerOperations {
	l := log.New("echo_test")
	return &GameServerOperations{
		Logger: l,
		Config: &ServiceConfig{
			Logger: l,
			Config: config.LoadAppConfig(),
		},
	}
}
func newWorldName(prefix string) string {
	return fmt.Sprintf("test-%s-%s", prefix, time.Now().Format(time.RFC3339))
}
