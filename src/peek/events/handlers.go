package events

import (
	"appengine"
	"appengine/user"
	"github.com/gorilla/mux"
	"net/http"
	"peek"
	"peek/ds"
	"strconv"
	"time"
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
	update, _ := strconv.ParseBool(r.FormValue("update"))

	c := appengine.NewContext(r)
	t := &uploadTask{
		context: c,
		events:  events,
		season:  season,
		league:  league,
		update:  update,
	}
	if err := t.exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/events/?s="+season, http.StatusFound)
}

const layout = "2006-01-02"

func index(w http.ResponseWriter, r *http.Request) {
	url := "/events/epl?s=2013-2014&d=" + time.Now().Format(layout)
	http.Redirect(w, r, url, http.StatusFound)
	return
}

func calc(w http.ResponseWriter, r *http.Request) {
	league := r.FormValue("league")
	season := r.FormValue("season")
	date, _ := time.Parse(layout, r.FormValue("date"))

	start, end := weekRange(date)
	c := appengine.NewContext(r)
	dst, keys, _ := ds.GetAllEventsByDateRange(c, league, season, start, end)

	events := make([]*Event, len(dst))
	for i, e := range dst {
		events[i] = &Event{
			Event: e,
			Id:    keys[i].IntID(),
		}
	}

	t := &calcTask{
		context: c,
		events:  events,
		season:  season,
		league:  league,
	}

	if err := t.exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := "/events/" + league + "?s=" + season + "&d=" + date.Format(layout)
	http.Redirect(w, r, url, http.StatusFound)
}

func resetView(w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/reset.html")
}

func reset(w http.ResponseWriter, r *http.Request) {
	league := r.FormValue("league")
	season := r.FormValue("season")
	t := &resetTask{
		context: appengine.NewContext(r),
		season:  season,
		league:  league,
	}

	if err := t.exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := "/events/" + league + "?s=" + season + "&d=" + time.Now().Format(layout)
	http.Redirect(w, r, url, http.StatusFound)
}

func league(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := vars["league"]
	season := r.FormValue("s")
	d := r.FormValue("d")

	if season == "" || d == "" {
		url := "/events/" + league + "?s=2013-2014&d=" + time.Now().Format(layout)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	date, _ := time.Parse(layout, d)

	c := appengine.NewContext(r)
	start, end := weekRange(date)
	dst, keys, _ := ds.GetAllEventsByDateRange(c, league, season, start, end)

	events := make([]*Event, len(dst))
	for i, e := range dst {
		events[i] = &Event{
			Event: e,
			Id:    keys[i].IntID(),
		}
	}

	gw := &GameWeek{
		Events:      events,
		PreviousUrl: "/events/" + league + "?s=" + season + "&d=" + date.AddDate(0, 0, -7).Format(layout),
		NextUrl:     "/events/" + league + "?s=" + season + "&d=" + date.AddDate(0, 0, 7).Format(layout),
		League:      league,
		Season:      season,
		Date:        date,
		IsAdmin:     user.IsAdmin(c),
	}

	peek.RenderTemplate(w, gw, "templates/events.html")
	return
}

func weekRange(date time.Time) (start, end time.Time) {
	y, m, d := date.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	days := (-1 * (int(date.Weekday()) + 1)) % 7
	start = today.AddDate(0, 0, days)
	end = start.AddDate(0, 0, 6)
	return
}
