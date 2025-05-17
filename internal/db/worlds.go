package db

import "github.com/ajaxe/mc-manager/internal/models"

func Worlds() (d []*models.WorldItem, err error) {
	var fn dbValFunc = func() any { return &models.WorldItem{} }

	r, err := readAllCollection(fn, collectionWorlds)

	d = make([]*models.WorldItem, len(r))
	for i, v := range r {
		d[i] = v.(*models.WorldItem)
	}

	return
}
