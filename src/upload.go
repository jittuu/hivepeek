package peek

import (
	"appengine"
	"appengine/datastore"
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

	dst, err := ParseEvents(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)
	for _, e := range dst {
		if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Event", nil), e); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/events/", http.StatusFound)
}
