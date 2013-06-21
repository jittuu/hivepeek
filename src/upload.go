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

	c := appengine.NewContext(r)
	createTeamsAndEvents(c, events)

	http.Redirect(w, r, "/events/", http.StatusFound)
}

func createTeamsAndEvents(c appengine.Context, events []*ds.Event) {
	teams, _, _ := ds.GetAllTeams(c)

	teamExists := func(name string) bool {
		for _, t := range teams {
			if t.Name == name {
				return true
			}
		}

		return false
	}

	visited := make(map[string]bool)

	for _, e := range events {
		if !visited[e.Home] && !teamExists(e.Home) {
			visited[e.Home] = true
			datastore.Put(c, datastore.NewIncompleteKey(c, "Team", nil), &ds.Team{Name: e.Home, Rating: 1000, RatingHome: 1000, RatingAway: 1000})
		}

		if !visited[e.Away] && !teamExists(e.Away) {
			visited[e.Away] = true
			datastore.Put(c, datastore.NewIncompleteKey(c, "Team", nil), &ds.Team{Name: e.Away, Rating: 1000, RatingHome: 1000, RatingAway: 1000})
		}

		datastore.Put(c, datastore.NewIncompleteKey(c, "Event", nil), e)
	}
}
