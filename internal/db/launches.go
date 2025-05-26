package db

import (
	"github.com/ajaxe/mc-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Launches(pgOpts PaginationOptions) (res models.PaginationResult[models.LaunchItem], err error) {
	var fn dbValFunc = func() any { return &models.LaunchItem{} }

	// default sort order is "descending" by "launch_date"
	sortOrder := -1
	if pgOpts.Direction == models.PageDirectionPrev {
		sortOrder = 1
	}

	cursor, err := pgOpts.Cursor()
	if err != nil {
		return
	}

	f := bson.D{}
	if pgOpts.CursorID != "" {
		if pgOpts.Direction == models.PageDirectionNext {
			f = append(f, bson.E{
				Key: "$or",
				Value: bson.A{
					bson.M{"launch_date": bson.M{"$lt": cursor.LaunchDate}},
					bson.D{
						{"launch_date", cursor.LaunchDate},
						{"_id", bson.M{"$lt": cursor.ID}},
					},
				},
			})
		} else {
			f = append(f, bson.E{
				Key: "$or",
				Value: bson.A{
					bson.M{"launch_date": bson.M{"$gt": cursor.LaunchDate}},
					bson.D{
						{"launch_date", cursor.LaunchDate},
						{"_id", bson.M{"$gt": cursor.ID}},
					},
				},
			})
		}
	}

	r, err := readAllCollection(readOptions{
		dbVal:      fn,
		collection: collectionLaunches,
		filter:     &f,
		opts: options.Find().
			SetSort(bson.D{
				{"launch_date", sortOrder},
				{"_id", sortOrder},
			}).
			SetLimit(int64(pgOpts.PageSize)),
	})

	d := make([]*models.LaunchItem, len(r))
	for i, v := range r {
		d[i] = v.(*models.LaunchItem)
	}

	if pgOpts.Direction == models.PageDirectionPrev {
		for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
			d[i], d[j] = d[j], d[i]
		}
	}

	count, err := collectionCount(collectionLaunches)

	prevID := ""
	if (pgOpts.Direction == models.PageDirectionPrev && len(d) == 0) || pgOpts.CursorID == "" {
		prevID = ""
	} else if len(d) != 0 {
		prevID = EncodePaginationCursor(&PaginationCursor{
			ID:         d[0].ID,
			LaunchDate: d[0].LaunchDate,
		})
	}

	nextID := ""
	if pgOpts.Direction == models.PageDirectionNext && len(d) == 0 {
		nextID = ""
	} else if len(d) != 0 {
		nextID = EncodePaginationCursor(&PaginationCursor{
			ID:         d[len(d)-1].ID,
			LaunchDate: d[len(d)-1].LaunchDate,
		})
	}

	res = models.PaginationResult[models.LaunchItem]{
		Total:   count,
		Results: d,
		PrevID:  prevID,
		NextID:  nextID,
	}

	return
}
func LaunchInsert(l *models.LaunchItem) (id bson.ObjectID, err error) {
	id = bson.NewObjectID()
	l.ID = id.Hex()

	err = insertRecord(l, collectionLaunches)
	return
}
