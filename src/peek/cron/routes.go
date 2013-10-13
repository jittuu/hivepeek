package cron

import (
	"github.com/mjibson/appstats"
	"peek"
)

func init() {
	r := peek.Router.PathPrefix("/cron").Subrouter()

	r.Handle("/pull", appstats.NewHandler(pull)).Methods("GET")
	r.Handle("/calc", appstats.NewHandler(calc)).Methods("GET")
	r.Handle("/fetch", appstats.NewHandler(fetch)).Methods("GET")
}
