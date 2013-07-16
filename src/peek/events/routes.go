package events

import (
	"peek"
)

func init() {
	s := peek.Router.PathPrefix("/events").Subrouter()

	s.HandleFunc("/new", newUpload).Methods("GET")
}
