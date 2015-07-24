package internal

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(v)
	return err
}
