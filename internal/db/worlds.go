package db

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Worlds() (d []*models.WorldItem, err error) {
	var fn dbValFunc = func() any { return &models.WorldItem{} }

	r, err := readAllCollection(fn, collectionWorlds)

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
