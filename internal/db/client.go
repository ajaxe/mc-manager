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
	collectionWorlds    = "worlds"
	collectionLaunches  = "launches"
	collectionPlaytimer = "playtimers"
)

const (
	readTimeout  = 30 * time.Second
	writeTimeout = 30 * time.Second
)

type Client struct {
	cli    *mongo.Client
	dbName string
}

func NewClient(c config.AppConfig) (*Client, error) {
	sAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(c.Database.ConnectionURI).
		SetServerAPIOptions(sAPI)

	opts.SetBSONOptions(&options.BSONOptions{
		ObjectIDAsHexString: true,
	})

	client, err := mongo.Connect(opts)

	if err != nil {
		return nil, err
	}

	if err = pingClient(client); err != nil {
		return nil, fmt.Errorf("ping faild to return in 2sec timeout: %v", err)
	}

	return &Client{
		cli:    client,
		dbName: c.Database.DbName,
	}, nil
}

func (c *Client) Ping() error {
	return pingClient(c.cli)
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

func (c *Client) Close(ctx context.Context) error {
	if c.cli == nil {
		log.Print("db client not instantiated, nothing to disconnect")
		return nil
	}
	return c.cli.Disconnect(ctx)
}
