package events

import (
	"appengine"
	"appengine/datastore"
	"peek/ds"
)

type resetTask struct {
	context appengine.Context
	season  string
	league  string
}

func (t *resetTask) exec() error {
	events, keys, err := ds.GetAllEvents(t.context, t.league, t.season)

	if err != nil {
		return err
	}

	visited := make(map[string]bool)

	for i, e := range events {
		if !visited[e.Home] {
			if err = t.resetTeam(e.HomeId); err != nil {
				return err
			}
		}

		if !visited[e.Away] {
			if err = t.resetTeam(e.AwayId); err != nil {
				return err
			}
		}

		if err = t.resetEvent(e, keys[i]); err != nil {
			return err
		}
	}

	return nil
}

func (t *resetTask) resetTeam(teamId int64) error {
	team := &ds.Team{}
	key := datastore.NewKey(t.context, "Team", "", teamId, nil)

	err := datastore.Get(t.context, key, team)
	if err != nil {
		return err
	}

	team.OverallRating = 1000
	team.OverallRatingLen = 0
	team.HomeNetRating = 0
	team.HomeNetRatingLen = 0
	team.AwayNetRating = 0
	team.AwayNetRatingLen = 0

	_, err = datastore.Put(t.context, key, team)
	return err
}

func (t *resetTask) resetEvent(e *ds.Event, key *datastore.Key) error {
	e.HRating = 0
	e.HRatingLen = 0
	e.HNetRating = 0
	e.HNetRatingLen = 0
	e.HFormRating = 0
	e.HFormRatingLen = 0

	e.ARating = 0
	e.ARatingLen = 0
	e.ANetRating = 0
	e.ANetRatingLen = 0
	e.AFormRating = 0
	e.AFormRatingLen = 0

	_, err := datastore.Put(t.context, key, e)

	return err
}
