package ds

import (
	"appengine"
	"appengine/datastore"
)

func GetTeam(c appengine.Context, name string, season string) (*Team, *datastore.Key, error) {
	q := datastore.
		NewQuery("Team").
		Filter("Name =", name).
		Filter("Season =", season).
		Limit(1)

	dst := make([]*Team, 0, 1)
	if keys, err := q.GetAll(c, &dst); err != nil {
		return nil, nil, err
	} else {
		if len(dst) > 0 {
			return dst[0], keys[0], nil
		} else {
			return nil, nil, nil
		}
	}
}

func GetAllEvents(c appengine.Context, league string, season string) ([]*Event, []*datastore.Key, error) {
	q := datastore.
		NewQuery("Event").
		Filter("Season =", season).
		Filter("League =", league)
	dst := make([]*Event, 0)

	if keys, err := q.GetAll(c, &dst); err != nil {
		return nil, nil, err
	} else {
		return dst, keys, nil
	}
}
