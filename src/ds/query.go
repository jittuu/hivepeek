package ds

import (
	"appengine"
	"appengine/datastore"
)

func GetAllTeams(c appengine.Context) ([]*Team, []*datastore.Key, error) {
	q := datastore.NewQuery("Team")
	dst := make([]*Team, 0)

	if keys, err := q.GetAll(c, &dst); err != nil {
		return nil, nil, err
	} else {
		return dst, keys, nil
	}
}

func GetTeam(c appengine.Context, name string) (*Team, *datastore.Key, error) {
	q := datastore.
		NewQuery("Team").
		Filter("Name =", name).
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

func GetAllEvents(c appengine.Context, season string) ([]*Event, []*datastore.Key, error) {
	q := datastore.NewQuery("Event").Filter("Season =", season)
	dst := make([]*Event, 0)

	if keys, err := q.GetAll(c, &dst); err != nil {
		return nil, nil, err
	} else {
		return dst, keys, nil
	}
}
