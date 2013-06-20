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

func GetAllEvents(c appengine.Context) ([]*Event, []*datastore.Key, error) {
	q := datastore.NewQuery("Event")
	dst := make([]*Event, 0)

	if keys, err := q.GetAll(c, &dst); err != nil {
		return nil, nil, err
	} else {
		return dst, keys, nil
	}
}
