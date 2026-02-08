package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type dbValFunc func() any

type readOptions struct {
	filter     *bson.D
	dbVal      dbValFunc
	opts       *options.FindOptionsBuilder
	collection string
}

func (c *Client) readAllCollection(ctx context.Context, ro readOptions) (d []any, err error) {
	if ro.dbVal == nil {
		err = fmt.Errorf("'dbVal' is required")
		return
	}
	if ro.collection == "" {
		err = fmt.Errorf("'collection' name is required")
		return
	}
	if ro.filter == nil {
		ro.filter = &bson.D{}
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	cur, err := c.cli.Database(c.dbName).
		Collection(ro.collection).
		Find(ctx, ro.filter, ro.opts)

	if err != nil {
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		r := ro.dbVal()
		if err = cur.Decode(r); err != nil {
			return
		}
		d = append(d, r)
	}

	return
}

func (c *Client) collectionCount(ctx context.Context, name string) (count int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	count, err = c.cli.Database(c.dbName).
		Collection(name).
		EstimatedDocumentCount(ctx)

	return
}

func (c *Client) deleteByID(ctx context.Context, id bson.ObjectID, collection string) (err error) {
	f := bson.D{{"_id", id.Hex()}}

	res, err := c.cli.Database(c.dbName).
		Collection(collection).
		DeleteMany(ctx, f)

	if err == nil {
		_ = res.DeletedCount
	}

	return
}

func (c *Client) insertRecord(ctx context.Context, u any, collection string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()

	_, err = c.cli.Database(c.dbName).
		Collection(collection).
		InsertOne(ctx, u)

	return
}

func (c *Client) readByID(ctx context.Context, id bson.ObjectID, v dbValFunc, collection string) (d any, err error) {
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	f := bson.D{{"_id", id.Hex()}}

	res := c.cli.Database(c.dbName).
		Collection(collection).
		FindOne(ctx, f)

	if res.Err() != nil {
		return nil, res.Err()
	}

	resp := v()

	err = res.Decode(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
