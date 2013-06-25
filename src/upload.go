package peek

import (
	"appengine"
	"appengine/datastore"
	"ds"
	"net/http"
)

func New(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, nil, "templates/upload.html")
}

func Create(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := ParseEvents(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	season := r.FormValue("season")

	c := appengine.NewContext(r)
	createTeams(c, events, season)
	createEvents(c, events, season)

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
	if t, k, _ := ds.GetTeam(c, name); t == nil {
		datastore.Put(
			c,
			datastore.NewIncompleteKey(c, "Team", nil),
			&ds.Team{
				Name: name,
				Ratings: []ds.SeasonRating{
					ds.SeasonRating{Season: season, Overall: 1000, Home: 1000, Away: 1000}}})
	} else {
		for _, r := range t.Ratings {
			if r.Season == season {
				return
			}
		}

		t.Ratings = append(t.Ratings, ds.SeasonRating{
			Season: season, Overall: 1000, Home: 1000, Away: 1000})

		datastore.Put(c, k, t)
	}
}

func createEvents(c appengine.Context, events []*ds.Event, season string) {
	existings, _, _ := ds.GetAllEvents(c, season)

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
			e.Season = season
			datastore.Put(c, datastore.NewIncompleteKey(c, "Event", nil), e)
		}
	}
}
