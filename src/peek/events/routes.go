package events

import (
	"peek"
)

func init() {
	r := peek.Router.PathPrefix("/events").Subrouter()

	r.HandleFunc("/new", newUpload).Methods("GET")
  r.HandleFunc("/upload", upload).Methods("POST")
}
