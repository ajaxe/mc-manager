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

func readAllCollection(ro readOptions) (d []any, err error) {
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

	c, err := NewClient()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	cur, err := c.Database(clientInstance.DbName).
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

func collectionCount(name string) (count int64, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	count, err = c.Database(clientInstance.DbName).
		Collection(name).
		EstimatedDocumentCount(ctx)

	return
}

func deleteByID(id bson.ObjectID, collection string) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	f := bson.D{{"_id", id.Hex()}}

	res, err := c.Database(clientInstance.DbName).
		Collection(collection).
		DeleteMany(context.TODO(), f)

	if err == nil {
		_ = res.DeletedCount
	}

	return
}

func insertRecord(u any, collection string) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	_, err = c.Database(clientInstance.DbName).
		Collection(collection).
		InsertOne(ctx, u)

	return
}

func readByID(id bson.ObjectID, v dbValFunc, collection string) (d any, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	f := bson.D{{"_id", id.Hex()}}

	res := c.Database(clientInstance.DbName).
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
