package events

import (
	"appengine"
	"net/http"
	"peek"
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
	t := &uploadTask{
		context: c,
		events:  events,
		season:  season,
		league:  league,
	}
	t.exec()

	http.Redirect(w, r, "/events/?s="+season, http.StatusFound)
}
