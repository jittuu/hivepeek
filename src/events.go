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

func Index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if dst, keys, err := ds.GetAllEvents(c); err != nil {
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
