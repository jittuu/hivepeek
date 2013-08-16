package events

import (
	"peek"
)

func init() {
	r := peek.Router.PathPrefix("/events").Subrouter()

	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/new", newUpload).Methods("GET")
	r.HandleFunc("/upload", upload).Methods("POST")
	r.HandleFunc("/calc", calc).Methods("POST")
	r.HandleFunc("/reset", resetView).Methods("GET")
	r.HandleFunc("/reset", reset).Methods("POST")
	r.HandleFunc("/run", runView).Methods("GET")
	r.HandleFunc("/run", run).Methods("POST")
	r.HandleFunc("/{league}", league).Methods("GET")
}
