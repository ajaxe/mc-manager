package db

import (
	"context"

	"github.com/ajaxe/mc-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (c *Client) Worlds(ctx context.Context) (d []*models.WorldItem, err error) {
	var fn dbValFunc = func() any { return &models.WorldItem{} }

	r, err := c.readAllCollection(ctx, readOptions{
		dbVal:      fn,
		collection: collectionWorlds,
	})

	d = make([]*models.WorldItem, len(r))
	for i, v := range r {
		d[i] = v.(*models.WorldItem)
	}

	return
}
func (c *Client) InsertWorld(ctx context.Context, w *models.WorldItem) (id bson.ObjectID, err error) {
	id = bson.NewObjectID()
	w.ID = id.Hex()

	err = c.insertRecord(ctx, w, collectionWorlds)
	return
}
func (c *Client) DeleteWorldByID(ctx context.Context, id bson.ObjectID) error {
	return c.deleteByID(ctx, id, collectionWorlds)
}
func (c *Client) WorldById(ctx context.Context, id bson.ObjectID) (w *models.WorldItem, err error) {
	var fn dbValFunc = func() any { return &models.WorldItem{} }

	r, err := c.readByID(ctx, id, fn, collectionWorlds)
	if err != nil {
		return nil, err
	}

	w = r.(*models.WorldItem)

	return
}
func (c *Client) UpdateWorldByID(ctx context.Context, id bson.ObjectID, w *models.WorldItem) (err error) {
	filter := bson.D{{"_id", id.Hex()}}
	if w.ID != id.Hex() {
		w.ID = id.Hex()
	}
	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()

	res, err := c.cli.Database(c.dbName).
		Collection(collectionWorlds).
		ReplaceOne(ctx, filter, w)

	if err == nil {
		_ = res.MatchedCount
	}

	return
}
