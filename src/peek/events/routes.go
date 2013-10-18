package events

import (
	"github.com/mjibson/appstats"
	"peek"
)

func init() {
	r := peek.Router.PathPrefix("/events").Subrouter()

	r.Handle("/", appstats.NewHandler(index)).Methods("GET")
	r.Handle("/qstats", appstats.NewHandler(qstats)).Methods("GET")
	r.Handle("/pull", appstats.NewHandler(getPull)).Methods("GET")
	r.Handle("/pull", appstats.NewHandler(pull)).Methods("POST")
	r.Handle("/fetch", appstats.NewHandler(fetchView)).Methods("GET")
	r.Handle("/fetch", appstats.NewHandler(fetch)).Methods("POST")
	r.Handle("/calc", appstats.NewHandler(calcView)).Methods("GET")
	r.Handle("/calc", appstats.NewHandler(calc)).Methods("POST")
	r.Handle("/reset", appstats.NewHandler(resetView)).Methods("GET")
	r.Handle("/reset", appstats.NewHandler(reset)).Methods("POST")
	r.Handle("/run", appstats.NewHandler(runView)).Methods("GET")
	r.Handle("/run", appstats.NewHandler(run)).Methods("POST")
	r.Handle("/today", appstats.NewHandler(today)).Methods("GET")
	r.Handle("/fixture", appstats.NewHandler(fixture)).Methods("GET")
	r.Handle("/dl/{league}/{season}", appstats.NewHandler(download)).Methods("GET")
	r.Handle("/{league}", appstats.NewHandler(league)).Methods("GET")
}
