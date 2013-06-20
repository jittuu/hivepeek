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
	if f, _, err := r.FormFile("file"); err == nil {
		dst := ParseEvents(f)
		c := appengine.NewContext(r)
		for _, e := range dst {
			_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Event", nil), e)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/events/", http.StatusFound)
	} else {
		panic(err)
	}
}
