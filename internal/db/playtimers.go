package db

import (
	"context"
	"time"

	"github.com/ajaxe/mc-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ActivePlayTimer returns play timer item with is_active set to true.
func (c *Client) ActivePlayTimer() (p *models.PlayTimerItem, err error) {
	f := bson.D{{"is_active", true}}

	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	err = c.cli.Database(c.dbName).
		Collection(collectionPlaytimer).
		FindOne(ctx, f).
		Decode(&p)

	if mongo.ErrNoDocuments == err {
		err = nil // No active play timer found, return nil
		p = nil
		return
	}

	return
}

// InsertPlayTimer inserts new instance of play timer item into the database.
func (c *Client) InsertPlayTimer(p *models.PlayTimerItem) (id bson.ObjectID, err error) {
	id = bson.NewObjectID()
	p.ID = id.Hex()

	err = c.insertRecord(p, collectionPlaytimer)
	return
}

// UpdatePlayTimerByID updates an existing play timer item by its ID.
func (c *Client) UpdatePlayTimerByID(id bson.ObjectID, p *models.PlayTimerItem) (err error) {
	filter := bson.D{{"_id", id.Hex()}}
	if p.ID != id.Hex() {
		p.ID = id.Hex()
	}
	p.LastUpdateDate = time.Now().UTC().Format(time.RFC3339)
	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	res, err := c.cli.Database(c.dbName).
		Collection(collectionPlaytimer).
		ReplaceOne(ctx, filter, p)

	if err == nil {
		_ = res.MatchedCount
	}

	return
}
