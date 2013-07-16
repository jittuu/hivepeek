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

func (t *uploadTask) exec() {
	t.createTeams()
	t.createEvents()
}

func (t *uploadTask) createTeams() {
	visited := make(map[string]bool)

	for _, e := range t.events {
		if !visited[e.Home] {
			visited[e.Home] = true
			t.createOrUpdateTeam(e.Home)
		}

		if !visited[e.Away] {
			visited[e.Away] = true
			t.createOrUpdateTeam(e.Away)
		}
	}
}

func (t *uploadTask) createOrUpdateTeam(name string) {
	if team, _, _ := ds.GetTeam(t.context, name, t.season); team == nil {
		datastore.Put(
			t.context,
			datastore.NewIncompleteKey(t.context, "Team", nil),
			&ds.Team{
				Name:          name,
				Season:        t.season,
				OverallRating: 1000,
				HomeNetRating: 0,
				AwayNetRating: 0})
	}
}

func (t *uploadTask) createEvents() {
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
			h, hk, _ := ds.GetTeam(t.context, e.Home, t.season)
			a, ak, _ := ds.GetTeam(t.context, e.Away, t.season)

			if h != nil && a != nil {
				e.HomeId = hk.IntID()
				e.AwayId = ak.IntID()
				e.Season = t.season
				e.League = t.league
				datastore.Put(t.context, datastore.NewIncompleteKey(t.context, "Event", nil), e)
			}
		}
	}
}
