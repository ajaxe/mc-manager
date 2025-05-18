package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type dbValFunc func() any

func readAllCollection(v dbValFunc, collection string) (d []any, err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	cur, err := c.Database(clientInstance.DbName).
		Collection(collection).
		Find(ctx, bson.D{})

	if err != nil {
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		r := v()
		if err = cur.Decode(r); err != nil {
			return
		}
		d = append(d, r)
	}

	return
}

func deleteByID(id bson.ObjectID, collection string) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	f := bson.D{{"_id", id}}

	res, err := c.Database(clientInstance.DbName).
		Collection(collection).
		DeleteMany(context.TODO(), f)

	_ = res.DeletedCount

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

	f := bson.D{{"_id", id}}

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
