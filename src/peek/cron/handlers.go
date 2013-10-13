package cron

import (
	"appengine"
	"net/http"
	"peek/events"
)

var (
	leagues = []string{"epl", "serie-a", "bundesliga", "la-liga", "ligue-1"}
)

func pull(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	f := func(context appengine.Context, l, s string) {
		events.DelayPull.Call(context, l, s, false)
	}

	handle(c, w, r, f)
}

func calc(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	f := func(context appengine.Context, l, s string) {
		events.DelayCalc.Call(context, l, s)
	}

	handle(c, w, r, f)
}

func handle(c appengine.Context, w http.ResponseWriter, r *http.Request, f func(c appengine.Context, l, s string)) {
	season := "2013-2014"
	for _, l := range leagues {
		f(c, l, season)
	}

	http.Redirect(w, r, "/events/qstats", http.StatusFound)
	return
}
