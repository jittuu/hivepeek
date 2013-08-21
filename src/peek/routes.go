package peek

import (
	"appengine"
	"github.com/gorilla/mux"
	"github.com/mjibson/appstats"
	"net/http"
)

var (
	Router = mux.NewRouter()
)

func init() {
	r := Router
	r.Handle("/", appstats.NewHandler(home))

	http.Handle("/", r)
}

func home(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, nil, "templates/home.html")
}
