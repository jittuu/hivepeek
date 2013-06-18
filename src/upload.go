package peek

import (
	"fmt"
	"net/http"
)

func New(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/upload.html")
}

func Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "new upload")
}
