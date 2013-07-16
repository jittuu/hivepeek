package events

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"peek"
	"peek/ds"
)

func newUpload(w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/upload.html")
}

func upload(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := parseEvents(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	season := r.FormValue("season")
	league := r.FormValue("league")

	c := appengine.NewContext(r)
	createTeams(c, events, season)
	createEvents(c, events, league, season)

	http.Redirect(w, r, "/events/?s="+season, http.StatusFound)
}

func createTeams(c appengine.Context, events []*ds.Event, season string) {
	visited := make(map[string]bool)

	for _, e := range events {
		if !visited[e.Home] {
			visited[e.Home] = true
			createOrUpdateTeam(c, e.Home, season)
		}

		if !visited[e.Away] {
			visited[e.Away] = true
			createOrUpdateTeam(c, e.Away, season)
		}
	}
}

func createOrUpdateTeam(c appengine.Context, name string, season string) {
	if t, _, _ := ds.GetTeam(c, name, season); t == nil {
		datastore.Put(
			c,
			datastore.NewIncompleteKey(c, "Team", nil),
			&ds.Team{
				Name:           name,
				Season:         season,
				OverallRating:  1000,
				HomeNetRating: 0,
				AwayNetRating: 0})
	}
}

func createEvents(c appengine.Context, events []*ds.Event, league string, season string) {
	existings, _, _ := ds.GetAllEvents(c, league, season)

	eventExists := func(e *ds.Event) bool {
		for _, de := range existings {
			if de.Away == e.Away && de.Home == e.Home {
				return true
			}
		}

		return false
	}

	for _, e := range events {
		if !eventExists(e) {
			h, hk, _ := ds.GetTeam(c, e.Home, season)
			a, ak, _ := ds.GetTeam(c, e.Away, season)

			if h != nil && a != nil {
				e.HomeId = hk.IntID()
				e.AwayId = ak.IntID()
				e.Season = season
				e.League = league
				datastore.Put(c, datastore.NewIncompleteKey(c, "Event", nil), e)
			}
		}
	}
}
