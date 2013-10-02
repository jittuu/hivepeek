package events

import (
	"appengine"
	"appengine/delay"
	"appengine/memcache"
	"appengine/urlfetch"
	"fmt"
)

var DelayPull = delay.Func("pull-result", PullResult)

func PullResult(c appengine.Context, league string, season string, update bool) {
	url := getPullUrl(league, season)
	client := urlfetch.Client(c)
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}

	events, err := parseEvents(resp.Body)
	if err != nil {
		panic(err)
	}

	t := &uploadTask{
		context: c,
		events:  events,
		season:  season,
		league:  league,
		update:  update,
	}
	err = t.exec()
	if err != nil {
		panic(err)
	}

	err = memcache.Flush(c)
	if err != nil {
		panic(err)
	}
}

func getPullUrl(l string, s string) string {
	spath := s[2:4] + s[7:]

	var fname string
	switch l {
	case "epl":
		fname = "E0.csv"
	case "serie-a":
		fname = "I1.csv"
	case "bundesliga":
		fname = "D1.csv"
	case "la-liga":
		fname = "SP1.csv"
	case "ligue-1":
		fname = "F1.csv"
	}
	return fmt.Sprintf("http://www.football-data.co.uk/mmz4281/%s/%s", spath, fname)
}
