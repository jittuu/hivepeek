package peek

import (
	"appengine"
	"ds"
	"net/http"
)

type Event struct {
	*ds.Event
	Id int64
}

type Team struct {
	*ds.Team
	Id int64
}

func Index(w http.ResponseWriter, r *http.Request) {
	var s = r.FormValue("s")
	if s == "" {
		http.Redirect(w, r, "/events/?s=2012-2013", http.StatusFound)
		return
	}

	c := appengine.NewContext(r)

	if dst, keys, err := ds.GetAllEvents(c, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		events := make([]*Event, len(dst))
		for i, e := range dst {
			events[i] = &Event{
				Event: e,
				Id:    keys[i].IntID()}
		}
		renderTemplate(w, events, "templates/events.html")
	}
}
