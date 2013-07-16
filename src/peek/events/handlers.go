package events

import (
	"net/http"
	"peek"
)

func newUpload(w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/upload.html")
}
