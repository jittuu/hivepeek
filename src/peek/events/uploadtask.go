package events

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"fmt"
	"peek/ds"
)

type uploadTask struct {
	context appengine.Context
	events  []*ds.Event
	season  string
	league  string
	update  bool
}

func (t *uploadTask) exec() error {
	teams, err := t.getOrAddTeams()
	if err != nil {
		return err
	}
	if err := t.createEvents(teams); err != nil {
		return err
	}

	return nil
}

func (t *uploadTask) getOrAddTeams() (map[string]*Team, error) {
	visited := make(map[string]*Team)

	for _, e := range t.events {
		if visited[e.Home] == nil {
			if h, err := t.getOrAddTeam(e.Home); err != nil {
				return nil, err
			} else {
				visited[e.Home] = h
			}
		}

		if visited[e.Away] == nil {
			if a, err := t.getOrAddTeam(e.Away); err != nil {
				return nil, err
			} else {
				visited[e.Away] = a
			}
		}
	}

	return visited, nil
}

func (t *uploadTask) getOrAddTeam(name string) (*Team, error) {
	team, k, err := ds.GetTeam(t.context, name, t.season)
	if err != nil {
		return nil, err
	}

	if team == nil {
		team = &ds.Team{
			Name:          name,
			Season:        t.season,
			OverallRating: 1000,
			HomeNetRating: 0,
			AwayNetRating: 0}
		k, err = datastore.Put(
			t.context,
			datastore.NewIncompleteKey(t.context, "Team", nil),
			team)
		if err != nil {
			return nil, err
		}
	}

	return &Team{Team: team, Id: k.IntID()}, nil
}

func (t *uploadTask) createEvents(teams map[string]*Team) error {
	oldEvents, err := t.getExistingEvents()
	if err != nil {
		return err
	}

	waits := make([]<-chan error, len(t.events))
	for i, e := range t.events {
		old := oldEvents.Find(e.Home, e.Away)
		waits[i] = t.createOrUpdateEvent(old, e, teams)
	}

	errorCount := 0
	for _, ch := range waits {
		err = <-ch
		if err != nil {
			errorCount++
		}
	}

	if errorCount > 0 {
		return errors.New(fmt.Sprintf("there are %d errors when saving %d events", errorCount, len(t.events)))
	}

	return nil
}

func (t *uploadTask) createOrUpdateEvent(old *Event, e *ds.Event, teams map[string]*Team) <-chan error {
	ch := make(chan error, 1)
	if old == nil {
		h := teams[e.Home]
		a := teams[e.Away]

		if h != nil && a != nil {
			e.HomeId = h.Id
			e.AwayId = a.Id
			e.Season = t.season
			e.League = t.league
			go func() {
				_, err := datastore.Put(
					t.context,
					datastore.NewIncompleteKey(t.context, "Event", nil), e)
				ch <- err
			}()
		}
	} else if t.update {
		old.AvgOdds = e.AvgOdds
		old.MaxOdds = e.MaxOdds
		key := datastore.NewKey(t.context, "Event", "", old.Id, nil)
		go func() {
			_, err := datastore.Put(t.context, key, old.Event)
			ch <- err
		}()
	} else {
    // we can send now since it is buffered (1)
		ch <- nil
	}

	return ch
}

func (t *uploadTask) getExistingEvents() (Events, error) {
	events, keys, err := ds.GetAllEvents(t.context, t.league, t.season)
	if err != nil {
		return nil, err
	}

	result := make([]*Event, len(events))
	for i, e := range events {
		result[i] = &Event{
			Event: e,
			Id:    keys[i].IntID(),
		}
	}

	return result, nil
}
