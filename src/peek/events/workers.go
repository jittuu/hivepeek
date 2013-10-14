package events

import (
	"appengine"
	"appengine/delay"
	"appengine/memcache"
	"appengine/urlfetch"
	"fmt"
)

var (
	DelayPull = delay.Func("pull-result", PullResult)

	DelayCalc = delay.Func("calc-result", CalcResult)

	DelayFetch = delay.Func("fetch-fixture", FetchFixture)
)

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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PullResult(c appengine.Context, league string, season string, update bool) {
	url := getPullUrl(league, season)
	client := urlfetch.Client(c)

	c.Infof("[start] getting data from %s", url)
	resp, err := client.Get(url)
	checkErr(err)
	c.Infof("[done] getting data from %s", url)

	c.Infof("[start] parsing data")
	events, err := parseEvents(resp.Body)
	checkErr(err)
	c.Infof("[done] parsing data with %d events", len(events))

	t := &uploadTask{
		context: c,
		events:  events,
		season:  season,
		league:  league,
		update:  update,
	}
	c.Infof("[start] uploading events")
	err = t.exec()
	checkErr(err)
	c.Infof("[done] uploading events")

	err = memcache.Flush(c)
	checkErr(err)
}

func CalcResult(c appengine.Context, league string, season string) {
	c.Infof("[start] calculating events for league: %s (%s)", league, season)

	t := &calcTask{
		context: c,
		season:  season,
		league:  league,
	}

	err := t.exec()
	checkErr(err)
	c.Infof("[done] calculating events for league: %s (%s)", league, season)

	err = memcache.Flush(c)
	checkErr(err)
}

func FetchFixture(c appengine.Context, league string) {
	c.Infof("[start] fetching fixtures for league: %s", league)
	task := &fetchTask{
		context: c,
		league:  league,
	}
	err := task.exec()
	checkErr(err)
	c.Infof("[done] fetching fixtures for league: %s", league)
}
