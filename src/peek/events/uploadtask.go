package events

import (
	"appengine"
	"appengine/datastore"
	"peek/ds"
)

type uploadTask struct {
	context appengine.Context
	events  []*ds.Event
	season  string
	league  string
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
	existings, _, _ := ds.GetAllEvents(t.context, t.league, t.season)

	eventExists := func(e *ds.Event) bool {
		for _, de := range existings {
			if de.Away == e.Away && de.Home == e.Home {
				return true
			}
		}

		return false
	}

	for _, e := range t.events {
		if !eventExists(e) {
			h := teams[e.Home]
			a := teams[e.Away]

			if h != nil && a != nil {
				e.HomeId = h.Id
				e.AwayId = a.Id
				e.Season = t.season
				e.League = t.league
				_, err := datastore.Put(
					t.context,
					datastore.NewIncompleteKey(t.context, "Event", nil), e)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
