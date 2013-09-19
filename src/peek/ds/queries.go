package ds

import (
	"appengine"
	"appengine/datastore"
	"time"
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

func GetAllTeams(c appengine.Context, season string) ([]*Team, []*datastore.Key, error) {
	q := datastore.
		NewQuery("Team").
		Filter("Season =", season)
	dst := make([]*Team, 0)

	if keys, err := q.GetAll(c, &dst); err != nil {
		return nil, nil, err
	} else {
		return dst, keys, nil
	}
}

func GetTeams(c appengine.Context, ids []int64) (teams []*Team, keys []*datastore.Key, err error) {
	keys = make([]*datastore.Key, len(ids))
	for i, id := range ids {
		keys[i] = datastore.NewKey(c, "Team", "", id, nil)
	}

	teams = make([]*Team, len(ids))
	for i, _ := range teams {
		teams[i] = &Team{}
	}
	err = datastore.GetMulti(c, keys, teams)
	return
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

func GetAllEventsByDateRange(c appengine.Context, league string, season string, start time.Time, end time.Time) ([]*Event, []*datastore.Key, error) {
	q := datastore.
		NewQuery("Event").
		Filter("Season =", season).
		Filter("League =", league).
		Filter("StartTime >=", start).
		Filter("StartTime <=", end)
	dst := make([]*Event, 0)

	if keys, err := q.GetAll(c, &dst); err != nil {
		return nil, nil, err
	} else {
		return dst, keys, nil
	}
}

func GetAllTeamMappings(c appengine.Context) (mappings []*TeamMapping, keys []*datastore.Key, err error) {
	q := datastore.NewQuery("TeamMapping")

	mappings = make([]*TeamMapping, 0)
	keys, err = q.GetAll(c, &mappings)
	return
}

func GetFixtures(c appengine.Context, start, end time.Time) (fixtures []*Fixture, keys []*datastore.Key, err error) {
	q := datastore.
		NewQuery("Fixture").
		Filter("StartTime >=", start).
		Filter("StartTime <=", end)

	fixtures = make([]*Fixture, 0)
	keys, err = q.GetAll(c, &fixtures)
	return
}

func GetFixturesByLeague(c appengine.Context, league string) (fixtures []*Fixture, keys []*datastore.Key, err error) {
	q := datastore.
		NewQuery("Fixture").
		Filter("League =", league)

	fixtures = make([]*Fixture, 0)
	keys, err = q.GetAll(c, &fixtures)
	return
}
