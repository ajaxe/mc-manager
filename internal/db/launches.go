package db

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Launches() (d []*models.LaunchItem, err error) {
	var fn dbValFunc = func() any { return &models.LaunchItem{} }

	r, err := readAllCollection(fn, collectionLaunches)

	d = make([]*models.LaunchItem, len(r))
	for i, v := range r {
		d[i] = v.(*models.LaunchItem)
	}

	return
}
func LaunchInsert(l *models.LaunchItem) (id bson.ObjectID, err error) {
	id = bson.NewObjectID()
	l.ID = id

	err = insertRecord(l, collectionLaunches)
	return
}
