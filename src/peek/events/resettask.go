package events

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"fmt"
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

	count := 0
	ch := make(chan error)
	for i, e := range events {
		if !visited[e.Home] {
			count++
			t.resetTeam(e.HomeId, ch)
		}

		if !visited[e.Away] {
			count++
			t.resetTeam(e.AwayId, ch)
		}

		count++
		t.resetEvent(e, keys[i], ch)
	}

	err_count := 0
	for i := 0; i < count; i++ {
		err = <-ch
		if err != nil {
			err_count++
		}
	}

	if err_count > 0 {
		return errors.New(fmt.Sprintf("there are %d errors when reseting", err_count))
	}

	return nil
}

func (t *resetTask) resetTeam(teamId int64, ch chan<- error) {
	team := &ds.Team{}
	key := datastore.NewKey(t.context, "Team", "", teamId, nil)

	err := datastore.Get(t.context, key, team)
	if err != nil {
		ch <- err
	}

	team.OverallRating = 1000
	team.OverallRatingLen = 0
	team.HomeNetRating = 0
	team.HomeNetRatingLen = 0
	team.AwayNetRating = 0
	team.AwayNetRatingLen = 0
	team.LastFiveMatchRating = make([]float64, 0)
	team.GoalsFor = ds.EventSector{}
	team.GoalsAgainst = ds.EventSector{}

	go func() {
		_, err = datastore.Put(t.context, key, team)
		ch <- err
	}()
}

func (t *resetTask) resetEvent(e *ds.Event, key *datastore.Key, ch chan<- error) {
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
	e.HGoalsFor = ds.EventSector{}
	e.HGoalsAgainst = ds.EventSector{}
	e.AGoalsFor = ds.EventSector{}
	e.AGoalsAgainst = ds.EventSector{}

	go func() {
		_, err := datastore.Put(t.context, key, e)
		ch <- err
	}()
}
