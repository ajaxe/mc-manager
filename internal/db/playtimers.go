package db

import (
	"context"
	"time"

	"github.com/ajaxe/mc-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ActivePlayTimer returns play timer item with is_active set to true.
func ActivePlayTimer() (p *models.PlayTimerItem, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	f := bson.D{{"is_active", true}}

	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	err = c.Database(clientInstance.DbName).
		Collection(collectionWorlds).
		FindOne(ctx, f).
		Decode(&p)

	return
}

// InsertPlayTimer inserts new instance of play timer item into the database.
func InsertPlayTimer(p *models.PlayTimerItem) (id bson.ObjectID, err error) {
	id = bson.NewObjectID()
	p.ID = id.Hex()

	err = insertRecord(p, collectionPlaytimer)
	return
}

// UpdatePlayTimerByID updates an existing play timer item by its ID.
func UpdatePlayTimerByID(id bson.ObjectID, p *models.PlayTimerItem) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	filter := bson.D{{"_id", id.Hex()}}
	if p.ID != id.Hex() {
		p.ID = id.Hex()
	}
	p.LastUpdateDate = time.Now().UTC().Format(time.RFC3339)
	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	res, err := c.Database(clientInstance.DbName).
		Collection(collectionPlaytimer).
		ReplaceOne(ctx, filter, p)

	if err == nil {
		_ = res.MatchedCount
	}

	return
}
