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
	r.Handle("/", internal.Handler(handler))

	lgs := r.PathPrefix("/leagues").Subrouter()
	lgs.Handle("/", internal.Handler(AllLeagues)).Methods("GET")
	lgs.Handle("/", internal.Handler(FetchLeagues)).Methods("POST")

	teams := r.PathPrefix("/teams").Subrouter()
	teams.Handle("/{league}/{season}", internal.Handler(AllTeamsByLeague)).Methods("GET")
	teams.Handle("/{league}/{season}", internal.Handler(FetchTeamsByLeague)).Methods("POST")

	mths := r.PathPrefix("/matches").Subrouter()
	mths.Handle("/{league}/{season}", internal.Handler(AllMatchesByLeague)).Methods("GET")
	mths.Handle("/{league}/{season}", internal.Handler(FetchMatchesByLeague)).Methods("POST")

	http.Handle("/", r)
}

func handler(c context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "Hello")
	return nil
}
