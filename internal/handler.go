package internal

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

type Handler func(context.Context, http.ResponseWriter, *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if err := h(c, w, r); err != nil {
		panic(err)
	}
}
