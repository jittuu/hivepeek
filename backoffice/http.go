package backoffice

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"github.com/jittuu/hivepeek/internal"
)

func init() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	cpt := r.PathPrefix("/leagues").Subrouter()
	cpt.Handle("/", internal.Handler(AllLeagues)).Methods("GET")
	cpt.Handle("/", internal.Handler(FetchLeagues)).Methods("POST")

	http.Handle("/", r)
}

func handler(c context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Hello, world!")
	return nil
}
