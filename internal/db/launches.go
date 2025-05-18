package db

import "github.com/ajaxe/mc-manager/internal/models"

func Launches() (d []*models.LaunchItem, err error) {
	var fn dbValFunc = func() any { return &models.LaunchItem{} }

	r, err := readAllCollection(fn, collectionLaunches)

	d = make([]*models.LaunchItem, len(r))
	for i, v := range r {
		d[i] = v.(*models.LaunchItem)
	}

	return
}
