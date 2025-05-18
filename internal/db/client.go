package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ajaxe/mc-manager/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	collectionWorlds   = "worlds"
	collectionLaunches = "launches"
)
const (
	readTimeout  = 30 * time.Second
	writeTimeout = 30 * time.Second
)

var clientInstance struct {
	Client *mongo.Client
	DbName string
}

func NewClient() (*mongo.Client, error) {
	cfg := config.LoadAppConfig()
	return NewClientWithConfig(cfg)
}

func NewClientWithConfig(c config.AppConfig) (*mongo.Client, error) {
	if clientInstance.Client != nil {
		return clientInstance.Client, nil
	}

	sAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(c.Database.ConnectionURI).
		SetServerAPIOptions(sAPI)

	client, err := mongo.Connect(opts)

	if err != nil {
		return nil, err
	}

	if err = pingClient(client); err != nil {
		return nil, fmt.Errorf("ping faild to return in 2sec timeout: %v", err)
	}

	clientInstance.Client = client
	clientInstance.DbName = c.Database.DbName

	return client, nil
}

func Ping() error {
	return pingClient(clientInstance.Client)
}
func pingClient(c *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if c == nil {
		return fmt.Errorf("client must be instantiated before calling Ping")
	}

	err := c.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("ping failed with 2sec timeout: %v", err)
	}
	return nil
}

func Terminate(ctx context.Context) error {
	client := clientInstance.Client
	if client == nil {
		log.Print("db client not instantiated, nothing to disconnect")
		return nil
	}
	e := client.Disconnect(ctx)
	client = nil
	return e
}
