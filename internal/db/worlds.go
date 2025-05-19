package db

import (
	"context"

	"github.com/ajaxe/mc-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Worlds() (d []*models.WorldItem, err error) {
	var fn dbValFunc = func() any { return &models.WorldItem{} }

	r, err := readAllCollection(readOptions{
		dbVal:      fn,
		collection: collectionWorlds,
	})

	d = make([]*models.WorldItem, len(r))
	for i, v := range r {
		d[i] = v.(*models.WorldItem)
	}

	return
}
func InsertWorld(w *models.WorldItem) (id bson.ObjectID, err error) {
	id = bson.NewObjectID()
	w.ID = id

	err = insertRecord(w, collectionWorlds)
	return
}
func DeleteWorldByID(id bson.ObjectID) error {
	return deleteByID(id, collectionWorlds)
}
func WorldById(id bson.ObjectID) (w *models.WorldItem, err error) {
	var fn dbValFunc = func() any { return &models.WorldItem{} }

	r, err := readByID(id, fn, collectionWorlds)
	if err != nil {
		return nil, err
	}

	w = r.(*models.WorldItem)

	return
}
func UpdateWorldByID(id bson.ObjectID, w *models.WorldItem) (err error) {
	c, err := NewClient()
	if err != nil {
		return
	}

	filter:= bson.D{{"_id", id }}
	w.ID = id
	ctx, cancel := context.WithTimeout(context.TODO(), writeTimeout)
	defer cancel()

	_, err = c.Database(clientInstance.DbName).
		Collection(collectionWorlds).
		ReplaceOne(ctx, filter, w)

	return
}
