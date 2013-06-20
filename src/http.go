package peek

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)

	s := r.PathPrefix("/events").Subrouter()

	s.HandleFunc("/", Index).Methods("GET")
	s.HandleFunc("/new", New).Methods("GET")
	s.HandleFunc("/create", Create).Methods("POST")

	http.Handle("/", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

const layout = "templates/layout.html"

func renderTemplate(w io.Writer, data interface{}, filenames ...string) {
	filenames = append(filenames, layout)
	if t, err := template.New("layout.html").ParseFiles(filenames...); err != nil {
		panic(err)
	} else {
		if tErr := t.Execute(w, data); tErr != nil {
			panic(tErr)
		}
	}
}
