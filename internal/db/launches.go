package db

import (
	"math"

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
			SetLimit(int64(pgOpts.PageSize * 2)),
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

	hasMore := len(d) > pgOpts.PageSize

	tot := len(d)

	if pgOpts.Direction == models.PageDirectionPrev {
		d = d[len(d)-pgOpts.PageSize:]
	} else {
		i := int(math.Min(float64(tot), float64(pgOpts.PageSize)))
		d = d[0:i]
	}

	prevID := EncodePaginationCursor(&PaginationCursor{
		ID:         d[0].ID,
		LaunchDate: d[0].LaunchDate,
	})
	if (pgOpts.Direction == models.PageDirectionPrev && !hasMore) || (pgOpts.Direction == models.PageDirectionNext && pgOpts.CursorID == "") {
		prevID = ""
	}

	nextID := EncodePaginationCursor(&PaginationCursor{
		ID:         d[len(d)-1].ID,
		LaunchDate: d[len(d)-1].LaunchDate,
	})
	if pgOpts.Direction == models.PageDirectionNext && !hasMore {
		nextID = ""
	}

	res = models.PaginationResult[models.LaunchItem]{
		Total:   count,
		Results: d,
		PrevID:  prevID,
		NextID:  nextID,
		HasMore: hasMore,
	}

	return
}
func LaunchInsert(l *models.LaunchItem) (id bson.ObjectID, err error) {
	id = bson.NewObjectID()
	l.ID = id.Hex()

	err = insertRecord(l, collectionLaunches)
	return
}
