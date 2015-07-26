package backoffice

import (
	"fmt"
	"net/http"
	"os"

	"google.golang.org/appengine/urlfetch"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"github.com/jittuu/hivepeek/internal"
	"github.com/jittuu/xmlsoccer"
)

func init() {
	r := mux.NewRouter()
	r.StrictSlash(true)
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
	fmt.Fprint(w, "Hello, world!")
	return nil
}

func Client(c context.Context) *xmlsoccer.Client {
	return &xmlsoccer.Client{
		BaseURL: xmlsoccer.DemoURL,
		APIKey:  os.Getenv("XMLSOCCER_API_KEY"),
		Client:  urlfetch.Client(c),
	}
}
