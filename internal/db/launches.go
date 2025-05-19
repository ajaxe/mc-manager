package db

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Launches() (d []*models.LaunchItem, err error) {
	var fn dbValFunc = func() any { return &models.LaunchItem{} }

	r, err := readAllCollection(readOptions{
		dbVal:      fn,
		collection: collectionLaunches,
		opts:       options.Find().SetSort(bson.D{{"launch_date", -1}}),
	})

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
