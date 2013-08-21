package events

import (
	"peek"
	"github.com/mjibson/appstats"
)

func init() {
	r := peek.Router.PathPrefix("/events").Subrouter()

	r.Handle("/", appstats.NewHandler(index)).Methods("GET")
	r.Handle("/new", appstats.NewHandler(newUpload)).Methods("GET")
	r.Handle("/upload", appstats.NewHandler(upload)).Methods("POST")
	r.Handle("/calc", appstats.NewHandler(calc)).Methods("POST")
	r.Handle("/reset", appstats.NewHandler(resetView)).Methods("GET")
	r.Handle("/reset", appstats.NewHandler(reset)).Methods("POST")
	r.Handle("/run", appstats.NewHandler(runView)).Methods("GET")
	r.Handle("/run", appstats.NewHandler(run)).Methods("POST")
	r.Handle("/{league}", appstats.NewHandler(league)).Methods("GET")
}
