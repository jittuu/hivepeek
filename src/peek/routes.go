package peek

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	Router = mux.NewRouter()
)

func init() {
	r := Router
	r.HandleFunc("/", home)


	http.Handle("/", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
